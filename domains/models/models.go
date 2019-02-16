package models

// These are domain model objects
// thread safe

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	Complete = iota
	Running  = iota
)

type Request struct {
	url         *url.URL
	Method      string
	Header      http.Header
	Body        string
	Cookie      []http.Cookie
	JobStatus   int
	NextRequest time.Time
	LastRequest time.Time
	Namespace   string

	status    map[string]float64
	semaphore chan struct{}
}

func NewRequest(
	namespace string,
	url *url.URL,
	method, body string) *Request {
	return &Request{
		Namespace: namespace,
		url:       url,
		Method:    method,
		Body:      body,
		status:    map[string]float64{},
		semaphore: make(chan struct{}, 1),
	}
}

func (r *Request) UrlString() string {
	r.semaphore <- struct{}{}
	defer func() {
		<-r.semaphore
	}()

	return r.url.String()
}

func (r *Request) SetStats(key string, value float64) {
	r.semaphore <- struct{}{}
	defer func() { <-r.semaphore }()

	if r.status == nil {
		r.status = make(map[string]float64)
	}
	r.status[key] = value
}

func (r *Request) GetStats(key string) float64 {
	r.semaphore <- struct{}{}
	defer func() { <-r.semaphore }()

	if r.status == nil {
		r.status = make(map[string]float64)
	}
	return r.status[key]
}

func (r *Request) GetStatsMap() map[string]float64 {
	return r.status
}

func (r *Request) SetStatusMap(m map[string]float64) {
	r.status = m
}

func (r *Request) CreateHTTPRequest() *http.Request {
	req, err := http.NewRequest(r.Method, r.UrlString(), strings.NewReader(r.Body))
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return req
}

func (r *Request) UrlHost() string {
	return r.url.Host
}

func (r *Request) UrlFragment() string {
	return r.url.Fragment
}

func (r *Request) UrlPath() string {
	return r.url.Path
}

func (r *Request) UrlScheme() string {
	return r.url.Scheme
}

type Response struct {
	Header     http.Header
	Body       []byte
	CreateAt   time.Time
	request    *http.Request
	Cookie     []http.Cookie
	StatusCode int
	Namespace  string
	semaphore  chan struct{}
}

func NewResponse(namespace string, resp *http.Response) *Response {
	r := &Response{}

	r.semaphore = make(chan struct{}, 1)
	r.Namespace = namespace
	r.request = resp.Request
	r.Cookie = make([]http.Cookie, len(resp.Cookies()))
	for _, c := range resp.Cookies() {
		r.Cookie = append(r.Cookie, *c)
	}
	r.StatusCode = resp.StatusCode
	r.CreateAt = time.Now()

	if resp.StatusCode == http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		} else {
			r.Body = b
			if err := resp.Body.Close(); err != nil {
				log.Println(err)
			}
		}
	}
	return r
}

func (r *Response) UrlString() string {
	r.semaphore <- struct{}{}
	defer func() { <-r.semaphore }()
	return r.request.URL.String()
}

func (r *Response) UrlParse(url string) (*url.URL, error) {
	r.semaphore <- struct{}{}
	defer func() { <-r.semaphore }()
	return r.request.URL.Parse(url)
}
