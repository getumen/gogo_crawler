package models

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	Complete = iota
	Running  = iota
)

type Request struct {
	url         *url.URL
	method      string
	header      http.Header
	body        []byte
	cookie      []*http.Cookie
	jobStatus   int
	nextRequest time.Time
	lastRequest time.Time
	stats       map[string]float64
	namespace   string
	mutex       sync.RWMutex
}

func NewRequest(Url *url.URL, Method string, Body []byte) *Request {
	return &Request{url: Url, method: Method, body: Body}
}

func NewRequestFromHTTP(request *http.Request) (*Request, error) {
	r := &Request{}
	r.url = request.URL
	r.method = request.Method
	r.header = request.Header
	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	err = request.Body.Close()
	if err != nil {
		return nil, err
	}
	r.body = b
	r.cookie = request.Cookies()
	r.stats = make(map[string]float64)
	r.mutex = sync.RWMutex{}
	return r, nil
}

func (r *Request) HttpRequest() (*http.Request, error) {
	req, err := http.NewRequest(
		r.method,
		r.url.String(),
		bytes.NewReader(r.Body()))
	if err != nil {
		return nil, err
	}
	for k, v := range r.header {
		for _, h := range v {
			req.Header.Add(k, h)
		}
	}
	for _, cookie := range r.cookie {
		req.AddCookie(cookie)
	}
	return req, nil
}

func (r *Request) UrlString() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.url.String()
}

func (r *Request) UrlHost() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.url.Host
}

func (r *Request) UrlPath() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.url.Path
}

func (r *Request) SetUrl(u *url.URL) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.url = u
}

func (r *Request) UrlParse(ref string) (*url.URL, error) {
	return r.url.Parse(ref)
}

func (r *Request) Method() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.method
}

func (r *Request) SetMethod(method string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.method = method
}

func (r *Request) Header() http.Header {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.header
}

func (r *Request) GetHeader(key string) string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.header.Get(key)
}

func (r *Request) AddHeader(key, value string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.header.Add(key, value)
}

func (r *Request) SetHeader(key, value string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.header.Set(key, value)
}

func (r *Request) DelHeader(key string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.header.Del(key)
}

func (r *Request) Body() []byte {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.body
}

func (r *Request) SetBody(body []byte) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.body = body
}

func (r *Request) AddCookie(cookie *http.Cookie) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.cookie = append(r.cookie, cookie)
}

func (r *Request) Cookies() []*http.Cookie {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.cookie
}

func (r *Request) SetCookies(cookies []*http.Cookie) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.cookie = cookies
}

func (r *Request) JobStatus() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.jobStatus
}

func (r *Request) SetJobStatus(jobStatus int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.jobStatus = jobStatus
}

func (r *Request) NextRequest() time.Time {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.nextRequest
}

func (r *Request) SetNextRequest(nextRequest time.Time) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.nextRequest = nextRequest
}

func (r *Request) LastRequest() time.Time {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.lastRequest
}

func (r *Request) SetLastRequest(lastRequest time.Time) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.lastRequest = lastRequest
}

func (r *Request) Stats(key string) float64 {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if r.stats == nil {
		r.stats = make(map[string]float64)
	}
	return r.stats[key]
}

func (r *Request) SetStats(key string, value float64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.stats == nil {
		r.stats = make(map[string]float64)
	}
	r.stats[key] = value
}

func (r *Request) StatsMap() map[string]float64 {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.stats
}

func (r *Request) SetStatsMap(m map[string]float64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.stats = m
}

func (r *Request) Namespace() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.namespace
}

func (r *Request) SetNamespace(namespace string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.namespace = namespace
}

type Response struct {
	header     http.Header
	body       []byte
	createAt   time.Time
	request    *http.Request
	cookie     []*http.Cookie
	statusCode int
	namespace  string
	mutex      sync.RWMutex
}

func NewResponseFromHTTP(response *http.Response) (*Response, error) {
	resp := &Response{}
	resp.header = response.Header
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}
	resp.body = b
	resp.request = response.Request
	resp.cookie = response.Cookies()
	resp.statusCode = response.StatusCode
	resp.mutex = sync.RWMutex{}
	return resp, nil
}

func (r *Response) Header() http.Header {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.header
}

func (r *Response) GetHeader(key string) string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.header.Get(key)
}

func (r *Response) AddHeader(key, value string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.header.Add(key, value)
}

func (r *Response) SetHeader(key, value string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.header.Set(key, value)
}

func (r *Response) DelHeader(key string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.header.Del(key)
}

func (r *Response) Body() []byte {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.body
}

func (r *Response) SetBody(body []byte) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.body = body
}

func (r *Response) CreateAt() time.Time {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.createAt
}

func (r *Response) SetCreateAt(createAt time.Time) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.createAt = createAt
}

func (r *Response) UrlString() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.request.URL.String()
}

func (r *Response) UrlHost() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.request.URL.Host
}

func (r *Response) UrlPath() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.request.URL.Path
}

func (r *Response) UrlParse(ref string) (*url.URL, error) {
	return r.request.URL.Parse(ref)
}

func (r *Response) Cookies() []*http.Cookie {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.cookie
}

func (r *Response) StatusCode() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.statusCode
}

func (r *Response) Namespace() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.namespace
}

func (r *Response) SetNamespace(namespace string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.namespace = namespace
}
