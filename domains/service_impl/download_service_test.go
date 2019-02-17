package service_impl

import (
	"github.com/getumen/gogo_crawler/domains/models"
	"github.com/getumen/gogo_crawler/domains/repository"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/url"
	"testing"
)

func TestConstructRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	clientMock := repository.NewMockHttpClientRepository(ctrl)
	s := downloadService{clientRepository: clientMock}

	u, _ := url.Parse("http://example.com")
	method := "GET"

	input := models.NewRequest(u, method, []byte{})
	expected, _ := http.NewRequest(method, u.String(), nil)

	actual, err := s.constructRequest(input)

	// TODO test body
	if err == nil &&
		expected.URL.String() == actual.URL.String() &&
		expected.Method == actual.Method {
	} else {
		t.Fatalf("%v != %v", expected, actual)
	}
}
