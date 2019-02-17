package service

import (
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
	"sync"
)

type DownloadService interface {
	DoRequest(ctx context.Context, in <-chan *models.Request, out chan<- *models.Response, wg *sync.WaitGroup)
}
