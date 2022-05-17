package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SseClient struct {
	messageChan chan string
	userName    string
}
type SseMessage struct {
	Message string `json:"message"`
	User    string `json:"user"`
}

func registerSseClient(w http.ResponseWriter, r *http.Request, hub *SseHub) {
	userName, err := getParam(r, "user")
	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, "The parameters have no names")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// For pushing data to clients, I call "flusher.Flush()"
	flusher, _ := w.(http.Flusher)

	newClient := SseClient{messageChan: make(chan string), userName: userName}
	hub.register <- &newClient

	defer func() {
		hub.unregister <- &newClient
	}()
	for {
		select {
		case message := <-newClient.messageChan:
			// push received messages to clients
			fmt.Fprintf(w, "data: %s\n\n", message)

			// Flush the data immediatly instead of buffering it for later.
			flusher.Flush()
		case <-r.Context().Done():
			// when "es.close()" is called, this loop operation will be ended.
			return
		}
	}
}
func sendSseMessage(w http.ResponseWriter, r *http.Request, hub *SseHub) {
	returnValue := &ActionResult{}
	w.Header().Set("Content-Type", "application/json")
	body, readBodyError := ioutil.ReadAll(r.Body)
	if readBodyError != nil {
		log.Println(readBodyError.Error())
		returnValue.Succeeded = false
		returnValue.ErrorMessage = "Failed reading values from body"
		failedReadingData, _ := json.Marshal(returnValue)
		w.Write(failedReadingData)
		return
	}
	message := &SseMessage{}
	convertError := json.Unmarshal(body, &message)
	if convertError != nil {
		log.Println(convertError.Error())
		returnValue.Succeeded = false
		returnValue.ErrorMessage = "Failed converting to WebSocketMessage"
		failedConvertingData, _ := json.Marshal(returnValue)
		w.Write(failedConvertingData)
		return
	}
	w.WriteHeader(200)
	hub.broadcast <- message.Message

	returnValue.Succeeded = true
	data, _ := json.Marshal(returnValue)
	w.Write(data)
}