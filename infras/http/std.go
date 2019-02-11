package http

import (
	"github.com/getumen/gogo_crawler/config"
	"github.com/getumen/gogo_crawler/domains/repository"
	"golang.org/x/net/proxy"
	"log"
	"net/http"
	"net/url"
)

type httpClientRepository http.Client

func NewHttpClientRepository(config *config.Config) repository.HttpClientRepository {
	if config.Http.ProxyUrl == "" {
		return &http.Client{}
	}
	u, err := url.Parse(config.Http.ProxyUrl)
	if err != nil {
		log.Fatal(err)
	}
	dialer, err := proxy.FromURL(u, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP transport
	// proxy does not support DialContext yet
	tr := &http.Transport{
		Dial: dialer.Dial,
	}
	return &http.Client{Transport: tr}
}
