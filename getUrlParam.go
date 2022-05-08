package main

import (
	"errors"
	"fmt"
	"net/http"
)

func getParam(r *http.Request, key string) (string, error) {
	result := r.URL.Query().Get(key)
	if len(result) <= 0 {
		return "", errors.New(fmt.Sprintf("no value: %s", key))
	}
	return result, nil
}
