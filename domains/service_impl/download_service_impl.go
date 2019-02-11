package service_impl

import (
	"context"
	"errors"
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/repository"
	"github.com/getumen/gogo_crawler/domains/service"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
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
			continue
		}
		resp, err := d.clientRepository.Do(req)
		if response, err := d.constructResponse(resp); err == nil {
			out <- response
		}
	}
}

func (d *downloadService) constructRequest(request *models.Request) (*http.Request, error) {
	r, err := http.NewRequest(request.Method, request.Url.String(), strings.NewReader(request.Body))
	if err != nil {
		return nil, err
	}
	for _, f := range d.requestMiddleware {
		f(r, request)
		if r == nil {
			return nil, errors.New("delete request")
		}
	}

	return r, nil
}

func (d *downloadService) constructResponse(response *http.Response) (*models.Response, error) {

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	r := &models.Response{}
	response.Request.URL.Fragment = ""

	r.Request = response.Request
	r.Cookie = make([]http.Cookie, len(response.Cookies()))
	for _, c := range response.Cookies() {
		r.Cookie = append(r.Cookie, *c)
	}
	r.StatusCode = response.StatusCode
	r.CreateAt = time.Now()

	if response.StatusCode == http.StatusOK {
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
		} else {
			r.Body = b
		}
	}

	for _, f := range d.responseMiddleware {
		f(response, r)
		if r == nil {
			return nil, errors.New("delete response")
		}
	}

	return r, nil
}
