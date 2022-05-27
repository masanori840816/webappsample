package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pion/webrtc/v3"
)

type SSEClient struct {
	userName string
	w        http.ResponseWriter
}

func newSSEClient(userName string, w http.ResponseWriter) *SSEClient {
	return &SSEClient{
		userName: userName,
		w:        w,
	}
}

func registerSSEClient(w http.ResponseWriter, r *http.Request, hub *SSEHub) {
	userName, err := getParam(r, "user")
	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, "The parameters have no names")
		return
	}
	newClient := newSSEClient(userName, w)
	peerConnection, err := NewPeerConnection()
	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, "Failed creating PeerConnection")
		return
	}
	ps, err := NewPeerConnectionState(newClient, peerConnection)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, "Failed connection")
		return
	}
	dc, err := NewWebRTCDataChannelStates(peerConnection)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, "Failed adding DataChannel")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	hub.register <- ps

	// For pushing data to clients, I call "flusher.Flush()"
	flusher, _ := w.(http.Flusher)
	defer func() {
		hub.unregister <- ps
		dc.Close()
		ps.Close()
	}()

	for {
		// handle PeerConnection events and close SSE event.
		select {
		case candidate := <-ps.candidateFound:
			jsonValue, err := NewCandidateMessageJSON(newClient.userName, candidate)
			if err != nil {
				log.Println(err.Error())
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", jsonValue)
			flusher.Flush()
		case track := <-ps.addTrack:
			hub.addTrack <- track
		case connectionState := <-ps.changeConnectionState:
			switch connectionState {
			case webrtc.PeerConnectionStateFailed:
				return
			case webrtc.PeerConnectionStateClosed:
				return
			}
		case message := <-dc.MessageCh:
			log.Println("GetMessage")
			log.Println(message.ID)
			log.Println(string(message.Message.Data))
		case <-r.Context().Done():
			// when "es.close()" is called, this loop operation will be ended.
			return
		}
	}
}
func sendSSEMessage(w http.ResponseWriter, r *http.Request, hub *SSEHub) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err.Error())
		j, _ := json.Marshal(GetFailed("Failed reading values from body"))

		w.Write(j)
		return
	}
	message := &ClientMessage{}
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Println(err.Error())
		j, _ := json.Marshal(GetFailed("Failed converting to ClientMessage"))

		w.Write(j)
		return
	}
	w.WriteHeader(200)
	hub.broadcast <- *message
	data, _ := json.Marshal(GetSucceeded())
	w.Write(data)
}
