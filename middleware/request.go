package middleware

import (
	"github.com/getumen/gogo_crawler/domains/models"
	"net/http"
)

var RequestMiddleware []func(httpRequest *http.Request, request *models.Request)

func init() {
	RequestMiddleware = append(RequestMiddleware, addHeader)
}

func addCookie(r *http.Request, model *models.Request) {
	for _, c := range model.Cookie {
		r.AddCookie(&c)
	}
}

func discardAllRequest(r *http.Request, model *models.Request) {
	r = nil
}

func addHeader(r *http.Request, model *models.Request) {
	headers := map[string]string{}
	headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	headers["Accept-Language"] = "ja-JP,ja;q=0.9,en-US;q=0.8,en;q=0.7"
	headers["Cache-Control"] = "max-age=0"
	headers["Host"] = model.Url.Host

	for k, v := range headers {
		r.Header.Set(k, v)
	}
}
