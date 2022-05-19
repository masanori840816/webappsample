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
	sendMessage           chan string
	candidateChan         chan *webrtc.ICECandidate
	changeConnectionState chan webrtc.PeerConnectionState
	addTrack              chan *webrtc.TrackRemote
	userName              string
}

func (c *SSEClient) CloseAllChannels() {
	close(c.sendMessage)
	close(c.candidateChan)
	close(c.changeConnectionState)
	close(c.addTrack)
}
func newSSEClient(userName string) *SSEClient {
	return &SSEClient{
		sendMessage:           make(chan string),
		candidateChan:         make(chan *webrtc.ICECandidate),
		changeConnectionState: make(chan webrtc.PeerConnectionState),
		addTrack:              make(chan *webrtc.TrackRemote),
		userName:              userName,
	}
}

func registerSSEClient(w http.ResponseWriter, r *http.Request, hub *SSEHub) {
	userName, err := getParam(r, "user")
	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, "The parameters have no names")
		return
	}
	newClient := newSSEClient(userName)
	ps, err := NewPeerConnectionState(newClient)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, "Failed connection")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// For pushing data to clients, I call "flusher.Flush()"
	flusher, _ := w.(http.Flusher)

	hub.register <- ps

	defer func() {
		hub.unregister <- ps
	}()
	for {
		select {
		case message := <-newClient.sendMessage:
			// push received messages to clients
			// This format must be like "data: {value}\n\n" or clients can't be gotten the value.
			fmt.Fprintf(w, "data: %s\n\n", message)

			// Flush the data immediatly instead of buffering it for later.
			flusher.Flush()
		case candidate := <-newClient.candidateChan:
			jsonValue, err := NewCandidateMessageJSON(newClient.userName, candidate)
			if err != nil {
				log.Println(err.Error())
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", jsonValue)
			flusher.Flush()
		case track := <-newClient.addTrack:
			hub.addTrack <- track
		case connectionState := <-newClient.changeConnectionState:
			switch connectionState {
			case webrtc.PeerConnectionStateFailed:
				if err := ps.peerConnection.Close(); err != nil {
					log.Print(err)
				}
			case webrtc.PeerConnectionStateClosed:
				signalPeerConnections(hub)
			}
		case <-r.Context().Done():
			// when "es.close()" is called, this loop operation will be ended.
			return
		}
	}
}
func sendSSEMessage(w http.ResponseWriter, r *http.Request, hub *SSEHub) {
	returnValue := &ActionResult{}
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err.Error())
		returnValue.Succeeded = false
		returnValue.ErrorMessage = "Failed reading values from body"
		failedReadingData, _ := json.Marshal(returnValue)
		w.Write(failedReadingData)
		return
	}
	message := &ClientMessage{}
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Println(err.Error())
		returnValue.Succeeded = false
		returnValue.ErrorMessage = "Failed converting to WebSocketMessage"
		failedConvertingData, _ := json.Marshal(returnValue)
		w.Write(failedConvertingData)
		return
	}
	w.WriteHeader(200)

	log.Println("sendMessage 4")
	hub.broadcast <- *message

	log.Println("sendMessage 5")
	returnValue.Succeeded = true
	data, _ := json.Marshal(returnValue)
	w.Write(data)

	log.Println("sendMessage 6")
}
