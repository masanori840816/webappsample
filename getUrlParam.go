package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetParam(r *http.Request, key string) (string, error) {
	result := r.URL.Query().Get(key)
	if len(result) <= 0 {
		return "", fmt.Errorf("no value: %s", key)
	}
	return result, nil
}

func GetClientMessage(w http.ResponseWriter, r *http.Request) (*ClientMessage, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed reading values from body: %s", err.Error())
		return nil, err
	}
	message := &ClientMessage{}
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Printf("Failed converting to ClientMessage: %s", err.Error())

		return nil, err
	}
	return message, err
}
