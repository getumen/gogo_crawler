package models

import "github.com/pkg/errors"

func ValidateResponse(r *Response) error {
	if r.namespace == "" {
		return errors.New("namespace is empty")
	}
	if r.statusCode == 0 {
		return errors.New("status code is invalid")
	}
	return nil
}
