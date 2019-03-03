package models

import (
	"errors"
	"net/url"
)

func ValidateRequest(r *Request) error {
	_, err := url.ParseRequestURI(r.url.String())
	if err != nil {
		return err
	}
	if r.namespace == "" {
		return errors.New("namespace is empty")
	}
	err = isValidMethod(r.method)
	if err != nil {
		return err
	}
	if r.stats == nil {
		return errors.New("stats map is nil")
	}
	return nil
}
