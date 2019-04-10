package models

import "errors"

func isValidMethod(method string) error {
	switch method {
	case "GET",
		"POST",
		"HEAD",
		"PUT":
		return nil
	}
	return errors.New("invalid method")
}
