package service

import (
	"context"
	"github.com/getumen/gogo_crawler/domains/models"
	"net/http"
	"sync"
)

type DownloadService interface {
	DoRequest(ctx context.Context, in <-chan *models.Request, out chan<- *models.Response, wg *sync.WaitGroup)
	AddRequestMiddleware(f func(r *http.Request, model *models.Request) *http.Request)
	AddResponseMiddleware(f func(r *http.Response, model *models.Response) *models.Response)
}
