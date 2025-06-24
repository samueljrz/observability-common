package util

import (
	"reflect"
)

func GetErrorName(err error) string {
	if err == nil {
		return ""
	}

	if errorType := reflect.TypeOf(err); errorType != nil && errorType.Name() != "" {
		return errorType.Name()
	}

	return "error"
}
