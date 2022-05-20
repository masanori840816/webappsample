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
	candidateChan         chan *webrtc.ICECandidate
	changeConnectionState chan webrtc.PeerConnectionState
	addTrack              chan *webrtc.TrackRemote
	userName              string
	w                     http.ResponseWriter
}

func (c *SSEClient) CloseAllChannels() {
	close(c.candidateChan)
	close(c.changeConnectionState)
	close(c.addTrack)
}
func newSSEClient(userName string, w http.ResponseWriter) *SSEClient {
	return &SSEClient{
		candidateChan:         make(chan *webrtc.ICECandidate),
		changeConnectionState: make(chan webrtc.PeerConnectionState),
		addTrack:              make(chan *webrtc.TrackRemote),
		userName:              userName,
		w:                     w,
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
	ps, err := NewPeerConnectionState(newClient)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, "Failed connection")
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
		ps.client.CloseAllChannels()
	}()
	for {
		select {
		case candidate := <-ps.client.candidateChan:
			jsonValue, err := NewCandidateMessageJSON(ps.client.userName, candidate)
			if err != nil {
				log.Println(err.Error())
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", jsonValue)
			flusher.Flush()
		case track := <-ps.client.addTrack:
			hub.addTrack <- track
		case connectionState := <-ps.client.changeConnectionState:
			switch connectionState {
			case webrtc.PeerConnectionStateFailed:
				break
			case webrtc.PeerConnectionStateClosed:
				break
			}
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
