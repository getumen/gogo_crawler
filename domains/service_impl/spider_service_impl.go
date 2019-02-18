package service_impl

import (
	"bytes"
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/service"
	"golang.org/x/net/html"
	"log"
	"regexp"
	"time"
)

type spiderService struct {
}

func NewSpiderService() service.SpiderService {
	return &spiderService{}
}

func (s *spiderService) ParseResponse(ctx context.Context, allowedDomainRegexp string, in <-chan *models.Response, out chan<- *models.Request) {
	allowedDomain := regexp.MustCompile(allowedDomainRegexp)
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
	close(out)
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
