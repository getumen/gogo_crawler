package usecase

import (
	"context"
	"github.com/getumen/gogo_crawler/config"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/service"
	"log"
	"sync"
)

type Crawler interface {
	Start(ctx context.Context, crawlerConfig *config.Config)
}

type crawler struct {
	crawlerService service.CrawlerService
}

func NewCrawler(crawlerService service.CrawlerService) Crawler {
	return &crawler{crawlerService: crawlerService}
}

func (c *crawler) Start(ctx context.Context, crawlerConfig *config.Config) {

	var wg sync.WaitGroup

	wg.Add(len(crawlerConfig.Page))
	for _, website := range crawlerConfig.Page {
		var website = &models.WebSite{
			Namespace:     website.Namespace,
			StartPage:     website.StartPage,
			AllowedDomain: website.AllowedDomain,
			DownloaderNum: website.DownloaderNum,
		}

		go func() {
			defer wg.Done()
			log.Printf("start crawl: %v\n", website)
			c.crawlerService.CrawlPage(ctx, website)
		}()
	}
	wg.Wait()
}
