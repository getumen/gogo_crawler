package service

import (
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
)

type CrawlerService interface {
	CrawlPage(ctx context.Context, website *models.WebSite)
}
