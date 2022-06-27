package main

import (
	"fmt"
	"net/http"
)

func getParam(r *http.Request, key string) (string, error) {
	result := r.URL.Query().Get(key)
	if len(result) <= 0 {
		return "", fmt.Errorf("no value: %s", key)
	}
	return result, nil
}
