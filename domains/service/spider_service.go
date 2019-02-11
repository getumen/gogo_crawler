package service

import (
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
)

type SpiderService interface {
	ParseResponse(ctx context.Context,
		allowedDomainRegexp string,
		in <-chan *models.Response,
		out chan<- *models.Request)
}
