package service_impl

import (
	"context"
	"errors"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/repository"
	"github.com/getumen/gogo_crawler/domains/service"
	"log"
	"net/http"
	"sync"
)

type downloadService struct {
	clientRepository   repository.HttpClientRepository
	requestMiddleware  []func(r *http.Request, model *models.Request)
	responseMiddleware []func(r *http.Response, model *models.Response)
}

func NewDownloadService(clientRepository repository.HttpClientRepository) service.DownloadService {
	return &downloadService{clientRepository: clientRepository}
}

func (d *downloadService) AddRequestMiddleware(f func(r *http.Request, model *models.Request)) {
	d.requestMiddleware = append(d.requestMiddleware, f)
}

func (d *downloadService) AddResponseMiddleware(f func(r *http.Response, model *models.Response)) {
	d.responseMiddleware = append(d.responseMiddleware, f)
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
		if response, err := d.constructResponse(request.Namespace, resp); err == nil {
			response.Namespace = request.Namespace
			out <- response
		} else {
			log.Printf("error constructResponse %v in %s", err, request.UrlString())
		}
	}
}

func (d *downloadService) constructRequest(request *models.Request) (*http.Request, error) {
	r := request.CreateHTTPRequest()
	for _, f := range d.requestMiddleware {
		f(r, request)
		if r == nil {
			return nil, errors.New("delete request")
		}
	}

	return r, nil
}

func (d *downloadService) constructResponse(namespace string, response *http.Response) (*models.Response, error) {

	r := models.NewResponse(namespace, response)
	for _, f := range d.responseMiddleware {
		f(response, r)
		if r == nil {
			return nil, errors.New("delete response")
		}
	}

	return r, nil
}
