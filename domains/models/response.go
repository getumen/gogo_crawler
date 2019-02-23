package models

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

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
