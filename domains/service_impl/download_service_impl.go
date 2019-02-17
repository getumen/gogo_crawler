package service_impl

import (
	"context"
	"errors"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/repository"
	"github.com/getumen/gogo_crawler/domains/service"
	"github.com/getumen/gogo_crawler/middleware"
	"log"
	"net/http"
	"sync"
	"time"
)

type downloadService struct {
	clientRepository repository.HttpClientRepository
}

func NewDownloadService(clientRepository repository.HttpClientRepository) service.DownloadService {
	return &downloadService{clientRepository: clientRepository}
}

func (d *downloadService) DoRequest(ctx context.Context, in <-chan *models.Request, out chan<- *models.Response, wg *sync.WaitGroup) {
	defer wg.Done()
	for request := range in {
		req, err := d.constructRequest(request)
		if err != nil {
			// request is deleted
			log.Printf("error constructRequest %v in %s", err, request.UrlString())
			continue
		}
		resp, err := d.clientRepository.Do(req)
		if err != nil {
			log.Printf("http request error %v in %s", err, request.UrlString())
			continue
		}
		if response, err := d.constructResponse(resp); err == nil {
			response.SetNamespace(request.Namespace())
			out <- response
		} else {
			log.Printf("error constructResponse %v in %s", err, request.UrlString())
		}
	}
}

func (d *downloadService) constructRequest(request *models.Request) (*http.Request, error) {
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

func (d *downloadService) constructResponse(response *http.Response) (*models.Response, error) {

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
