package usecase

import (
	"context"
	"github.com/getumen/gogo_crawler/config"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/service"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type Crawler interface {
	Start(ctx context.Context)
	RequestMiddleware(f func(r *http.Request, model *models.Request))
	ResponseMiddleware(f func(r *http.Response, model *models.Response))
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

func (c *crawler) RequestMiddleware(f func(r *http.Request, model *models.Request)) {
	c.downloadService.AddRequestMiddleware(f)
}

func (c *crawler) ResponseMiddleware(f func(r *http.Response, model *models.Response)) {
	c.downloadService.AddResponseMiddleware(f)
}

func (c *crawler) Start(ctx context.Context) {
	crawlerGroup := &sync.WaitGroup{}

	downloaderNum := c.conf.Crawler.DownloaderNum

	for _, website := range c.conf.Page {
		crawlerGroup.Add(1)
		log.Printf("start crawl: %v\n", website)
		go func(website config.WebSite) {

			defer crawlerGroup.Done()

			scheduledReqChan := make(chan *models.Request, downloaderNum)

			downloadRespChan := make(chan *models.Response, downloaderNum)

			wg := &sync.WaitGroup{}
			wg.Add(downloaderNum)
			for i := 0; i < downloaderNum; i++ {
				go c.downloadService.DoRequest(ctx, scheduledReqChan, downloadRespChan, wg)
			}
			go func() {
				wg.Wait()
				close(downloadRespChan)
			}()

			scheduledRespChan := make(chan *models.Response, downloaderNum)
			itemRespChan := make(chan *models.Response, downloaderNum)
			spiderRespChan := make(chan *models.Response, downloaderNum)

			go distributeResponse(downloadRespChan, scheduledRespChan, itemRespChan, spiderRespChan)

			go c.scheduleService.ScheduleRequest(ctx, scheduledRespChan)

			go c.itemService.SaveResponse(ctx, itemRespChan)

			scheduleReqChan := make(chan *models.Request, downloaderNum)
			go c.spiderService.ParseResponse(ctx, website.AllowedDomain, spiderRespChan, scheduleReqChan)

			go c.scheduleService.ScheduleNewRequest(ctx, scheduleReqChan)

			u, err := url.Parse(website.StartPage)
			if err == nil {
				scheduledReqChan <- models.NewRequest(website.Namespace, u, "GET", "")
			}

			c.scheduleService.GenerateRequest(ctx, website.Namespace, scheduledReqChan)
		}(website)
	}

	crawlerGroup.Wait()
}

func distributeResponse(in <-chan *models.Response, out ...chan<- *models.Response) {
	for resp := range in {
		for i := 0; i < len(out); i++ {
			out[i] <- resp
		}
	}
	for i := 0; i < len(out); i++ {
		close(out[i])
	}
}
