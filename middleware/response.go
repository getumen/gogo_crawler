package middleware

import (
	"github.com/getumen/gogo_crawler/domains/models"
	"net/http"
)

var ResponseMiddleware []func(httpResponse *http.Response, response *models.Response)

func init() {
	ResponseMiddleware = append(ResponseMiddleware, doNothing)
}

func doNothing(httpResponse *http.Response, response *models.Response) {
}
