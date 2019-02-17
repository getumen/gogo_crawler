package usecase

import (
	"context"
	"github.com/getumen/gogo_crawler/config"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/service"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Crawler interface {
	Start(ctx context.Context)
}

type crawler struct {
	conf            *config.Config
	scheduleService service.ScheduleService
	downloadService service.DownloadService
	spiderService   service.SpiderService
	itemService     service.ItemService
}

func NewCrawler(
	conf *config.Config,
	scheduleService service.ScheduleService,
	downloadService service.DownloadService,
	spiderService service.SpiderService,
	itemService service.ItemService) Crawler {
	return &crawler{
		conf:            conf,
		scheduleService: scheduleService,
		downloadService: downloadService,
		spiderService:   spiderService,
		itemService:     itemService,
	}
}

var signalChan = make(chan os.Signal, 1)

func (c *crawler) Start(ctx context.Context) {
	crawlerGroup := &sync.WaitGroup{}

	downloaderNum := c.conf.Crawler.DownloaderNum

	cancelCtx, cancel := context.WithCancel(ctx)

	go func() {
		cnt := 0
		signal.Notify(signalChan,
			syscall.SIGINT)
		for {
			s := <-signalChan
			switch s {
			case syscall.SIGINT:
				log.Println("SIGINT")
				cnt += 1
				if cnt > 1 {
					os.Exit(0)
				} else {
					cancel()
				}
			default:
				log.Fatalln("Unknown signal.")
			}
		}
	}()

	for _, website := range c.conf.Page {
		crawlerGroup.Add(1)
		log.Printf("start crawl: %v\n", website)
		go func(website config.WebSite) {

			sinkWaitGroup := &sync.WaitGroup{}

			scheduledReqChan := make(chan *models.Request)

			downloadRespChan := make(chan *models.Response)

			wg := &sync.WaitGroup{}
			wg.Add(downloaderNum)
			for i := 0; i < downloaderNum; i++ {
				go c.downloadService.DoRequest(cancelCtx, scheduledReqChan, downloadRespChan, wg)
			}
			go func() {
				wg.Wait()
				close(downloadRespChan)
			}()

			scheduledRespChan := make(chan *models.Response)
			itemRespChan := make(chan *models.Response)
			spiderRespChan := make(chan *models.Response)

			go distributeResponse(downloadRespChan, scheduledRespChan, itemRespChan, spiderRespChan)

			go func() {
				sinkWaitGroup.Add(1)
				c.scheduleService.ScheduleRequest(cancelCtx, scheduledRespChan)
				sinkWaitGroup.Done()
			}()
			go func() {
				sinkWaitGroup.Add(1)
				c.itemService.SaveResponse(cancelCtx, itemRespChan)
				sinkWaitGroup.Done()
			}()

			scheduleReqChan := make(chan *models.Request)
			go c.spiderService.ParseResponse(cancelCtx, website.AllowedDomain, spiderRespChan, scheduleReqChan)

			go func() {
				sinkWaitGroup.Add(1)
				c.scheduleService.ScheduleNewRequest(cancelCtx, scheduleReqChan)
				sinkWaitGroup.Done()
			}()
			u, err := url.Parse(website.StartPage)
			if err == nil {
				req := models.NewRequest(u, "GET", nil)
				req.SetNamespace(website.Namespace)
				scheduledReqChan <- req
			}

			go c.scheduleService.GenerateRequest(cancelCtx, website.Namespace, scheduledReqChan)

			sinkWaitGroup.Wait()
			crawlerGroup.Done()

		}(website)
	}

	crawlerGroup.Wait()
}

func distributeResponse(in <-chan *models.Response, out ...chan<- *models.Response) {
	for resp := range in {
		for i := 0; i < len(out); i++ {
			response := resp
			out[i] <- response
		}
	}
	for i := 0; i < len(out); i++ {
		close(out[i])
	}
}
