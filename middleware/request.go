package middleware

import (
	"github.com/getumen/gogo_crawler/domains/models"
	"net/http"
)

func AddCookie(r *http.Request, model *models.Request) {
	for _, c := range model.Cookie {
		r.AddCookie(&c)
	}
}

func DiscardAllRequest(r *http.Request, model *models.Request) {
	r = nil
}
