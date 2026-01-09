package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	messages "webappsample/messages"
)

func GetParam(r *http.Request, key string) (string, error) {
	result := r.URL.Query().Get(key)
	if len(result) <= 0 {
		return "", fmt.Errorf("no value: %s", key)
	}
	return result, nil
}

func GetClientMessage(w http.ResponseWriter, r *http.Request) (*messages.ClientMessage, error) {
	message := &messages.ClientMessage{}
	err := json.NewDecoder(r.Body).Decode(message)
	if err != nil {
		log.Printf("Failed converting to ClientMessage: %s", err.Error())

		return nil, err
	}
	return message, err
}
