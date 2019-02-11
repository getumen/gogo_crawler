package middleware

import (
	"github.com/getumen/gogo_crawler/domains/models"
	"net/http"
)

func AddCookie(r *http.Request, model *models.Request) *http.Request {
	for _, c := range model.Cookie {
		r.AddCookie(&c)
	}
	return r
}
