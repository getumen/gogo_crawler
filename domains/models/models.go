package models

import (
	"net/http"
	"net/url"
	"time"
)

const (
	Complete = iota
	Running  = iota
)

type Request struct {
	Url         *url.URL
	Method      string
	Header      http.Header
	Body        string
	Cookie      []http.Cookie
	JobStatus   int
	NextRequest time.Time
	LastRequest time.Time
	Stats       map[string]float64
	Namespace   string
}

type Response struct {
	Header     http.Header
	Body       []byte
	CreateAt   time.Time
	Request    *http.Request
	Cookie     []http.Cookie
	StatusCode int
	Namespace  string
}
