package service_impl

import (
	"bytes"
	"context"
	"errors"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/repository"
	"github.com/getumen/gogo_crawler/domains/service"
	"github.com/getumen/gogo_crawler/middleware"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"
)

const generateRequestLimit = 10
const heartBeat = time.Second
const channelSize = 10

type crawlerService struct {
	clientRepository   repository.HttpClientRepository
	requestRepository  repository.RequestRepository
	responseRepository repository.ResponseRepository
	scheduleRule       service.ScheduleRule
}

func NewCrawlerService(
	clientRepository repository.HttpClientRepository,
	requestRepository repository.RequestRepository,
	responseRepository repository.ResponseRepository,
	scheduleRule service.ScheduleRule, ) service.CrawlerService {
	return &crawlerService{clientRepository: clientRepository,
		requestRepository:  requestRepository,
		responseRepository: responseRepository,
		scheduleRule:       scheduleRule}
}

func (c *crawlerService) CrawlPage(ctx context.Context, website *models.WebSite) {

	requestPipeline := c.generateRequest(ctx, website.Namespace)

	downloadPipeline := c.doRequest(ctx, requestPipeline, website.DownloaderNum)

	respPipeline1, respPipeline2, respPipeline3 := tee(downloadPipeline)

	newRequestPipeline := c.parseResponse(ctx, website.AllowedDomain, respPipeline1)

	out1 := c.scheduleNewRequest(ctx, newRequestPipeline)
	out2 := c.scheduleRequest(ctx, respPipeline2)
	out3 := c.saveResponse(ctx, respPipeline3)

	out := funIn(out1, out2, out3)

	u, err := url.ParseRequestURI(website.StartPage)
	if err == nil {
		req := models.NewRequest(u, "GET", nil)
		req.SetNamespace(website.Namespace)
		req.SetLastRequest(time.Now())
		newRequestPipeline <- req
	}
	for range out {
	}
}

func (c *crawlerService) doRequest(ctx context.Context, in <-chan *models.Request, downloaderNum int) chan *models.Response {

	out := make(chan *models.Response, channelSize)
	var wg sync.WaitGroup

	wg.Add(downloaderNum)

	go func() {
		wg.Wait()
		close(out)
	}()

	for i := 0; i < downloaderNum; i++ {
		go func() {
			defer wg.Done()
			for request := range in {
				req, err := c.constructRequest(request)
				if err != nil {
					// request is deleted
					log.Printf("error constructRequest %v in %s", err, request.UrlString())
					continue
				}
				resp, err := c.clientRepository.Do(req)
				if err != nil {
					log.Printf("http request error %v in %s", err, request.UrlString())
					continue
				}
				if response, err := c.constructResponse(resp); err == nil {
					response.SetNamespace(request.Namespace())
					out <- response
				} else {
					log.Printf("error constructResponse %v in %s", err, request.UrlString())
				}
			}

		}()
	}
	return out
}

func (c *crawlerService) constructRequest(request *models.Request) (*http.Request, error) {
	r, err := request.HttpRequest()
	if err != nil {
		return nil, err
	}
	for _, f := range middleware.RequestMiddleware {
		f(r, request)
		if r == nil {
			return nil, errors.New("delete request")
		}
	}

	return r, nil
}

func (c *crawlerService) constructResponse(response *http.Response) (*models.Response, error) {

	r, err := models.NewResponseFromHTTP(response)
	if err != nil {
		return nil, err
	}
	r.SetCreateAt(time.Now())

	for _, f := range middleware.ResponseMiddleware {
		f(response, r)
		if r == nil {
			return nil, errors.New("delete response")
		}
	}

	return r, nil
}

func (c *crawlerService) saveResponse(ctx context.Context, in <-chan *models.Response) chan interface{} {
	out := make(chan interface{}, channelSize)
	go func() {
		defer close(out)
		for response := range in {
			if exist, err := c.responseRepository.IsExist(ctx, response); !exist && err == nil {
				_ = c.responseRepository.Save(ctx, response)
			} else if err != nil {
				log.Println(err)
			}
		}
	}()
	return out
}

func (c *crawlerService) generateRequest(ctx context.Context, namespace string) chan *models.Request {

	out := make(chan *models.Request, channelSize)

	go func() {
		defer func() {
			log.Printf("shutdown GenerateRequest")
			close(out)
		}()

		ticker := time.NewTicker(heartBeat)
		for {
			select {
			case <-ticker.C:
				requests, err := c.requestRepository.FindAllByDomainAndBeforeTimeOrderByNextRequest(
					ctx,
					namespace,
					time.Now(),
					0,
					generateRequestLimit)

				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("Scheduling %d request in %s\n", len(requests), namespace)
				for _, request := range requests {
					out <- request
				}
			case <-ctx.Done():
				return
			default:
				time.Sleep(heartBeat)
			}
		}
	}()

	return out
}

func (c *crawlerService) scheduleRequest(ctx context.Context, in <-chan *models.Response) chan interface{} {
	out := make(chan interface{}, channelSize)

	go func() {

		defer close(out)

		for resp := range in {
			r, err := c.requestRepository.FindByUrl(ctx, resp.Namespace(), resp.UrlString())
			if err != nil {
				log.Println(err)
				continue
			}

			if resp.StatusCode() == http.StatusOK {

				exist, err := c.responseRepository.IsExist(ctx, resp)
				if err != nil {
					log.Println(err)
					continue
				}

				if exist {
					log.Printf("Not updated url: %s\n", r.UrlString())
					r.SetNextRequest(time.Now().Add(time.Since(r.LastRequest()) * 2))
				} else {
					log.Printf("updated url: %s\n", r.UrlString())

					c.scheduleRule.UpdateStatsSuccess(r)

					r.SetLastRequest(time.Now())

					// schedule next request by schedule rule
					r.SetNextRequest(c.scheduleRule.ScheduleNextSuccess(r))
				}
			} else {
				log.Printf("http error %d in %s", resp.StatusCode(), resp.UrlString())
				c.scheduleRule.UpdateStatsFail(r)
				r.SetNextRequest(c.scheduleRule.ScheduleNextFail(r))
			}

			r.SetNamespace(resp.Namespace())
			err = c.requestRepository.Save(ctx, r)
			if err != nil {
				log.Println(err)
				continue
			}
		}

	}()

	return out
}

func (c *crawlerService) scheduleNewRequest(ctx context.Context, in <-chan *models.Request) chan interface{} {
	out := make(chan interface{}, channelSize)
	go func() {
		defer close(out)
		for request := range in {
			exist, err := c.requestRepository.IsExist(ctx, request.Namespace(), request.UrlString())
			if err != nil {
				log.Println(err)
				continue
			}
			if !exist {
				request.SetNextRequest(request.LastRequest().Add(time.Duration(1) * time.Minute))
				request.SetStats("AccessCount", 1)
				request.SetStats("AccessSuccess", 1)
				request.SetStats("IntervalSum", 0)
				err = c.requestRepository.Save(ctx, request)
				if err != nil {
					log.Println(err)
					continue
				} else {
					log.Printf("Schedule new request: %s\n", request.UrlString())
				}
			}
		}
	}()

	return out
}

func (c *crawlerService) parseResponse(ctx context.Context, allowedDomainRegexp string, in <-chan *models.Response) chan *models.Request {

	out := make(chan *models.Request, channelSize)

	allowedDomain := regexp.MustCompile(allowedDomainRegexp)

	go func() {
		defer close(out)
		for response := range in {
			visitNode := func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "a" {
					for _, a := range n.Attr {
						if a.Key != "href" {
							continue
						}
						link, err := response.UrlParse(a.Val)
						if err != nil {
							log.Printf("html parse error\n")
							continue // ignore bad URLs
						}
						if !allowedDomain.MatchString(link.Host) {
							continue // ignore not allowed domain
						}
						link.Fragment = ""
						req := models.NewRequest(link, "GET", nil)
						req.SetCookies(response.Cookies())
						req.SetLastRequest(time.Now())
						req.SetNamespace(response.Namespace())
						out <- req
					}
				}
			}
			buf := bytes.NewBuffer(response.Body())
			doc, err := html.Parse(buf)
			if err != nil {
				log.Println(err)
			}
			forEachNode(doc, visitNode, nil)
		}
	}()

	return out
}

func tee(in <-chan *models.Response) (_, _, _ chan *models.Response) {
	const n = 3
	out1 := make(chan *models.Response, channelSize)
	out2 := make(chan *models.Response, channelSize)
	out3 := make(chan *models.Response, channelSize)

	go func() {
		defer func() {
			close(out1)
			close(out2)
			close(out3)
		}()

		for resp := range in {
			var out1, out2, out3 = out1, out2, out3
			for i := 0; i < n; i++ {
				select {
				case out1 <- resp:
					out1 = nil
				case out2 <- resp:
					out2 = nil
				case out3 <- resp:
					out3 = nil
				}
			}
		}
	}()
	return out1, out2, out3
}

func funIn(channels ...<-chan interface{}) chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			multiplexedStream <- i
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
