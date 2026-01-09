package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pion/webrtc/v4"

	messages "webappsample/messages"
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
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(200)
	// For pushing data to clients, I call "flusher.Flush()"
	flusher, _ := w.(http.Flusher)
	userName, err := GetParam(r, "user")
	if err != nil {
		log.Println(err.Error())
		rejectConnection(w, r, "Missing or invalid user parameter")
		return
	}
	// Duplicated username error
	for c := range hub.clients {
		if c.client.userName == userName {
			rejectConnection(w, r, "The user name is already used")
			return
		}
	}
	newClient := newSSEClient(userName, w)
	peerConnection, err := NewPeerConnection()
	if err != nil {
		log.Println(err.Error())
		rejectConnection(w, r, "Failed to connect Server Sent Events")
		return
	}
	dc, err := NewWebRTCDataChannelStates(peerConnection)
	if err != nil {
		peerConnection.Close()
		log.Println(err.Error())
		rejectConnection(w, r, "Failed to create data channel")
		return
	}
	ps, err := NewPeerConnectionState(newClient, peerConnection, dc)
	if err != nil {
		log.Println(err.Error())
		rejectConnection(w, r, "Failed to connect Server Sent Events")
		return
	}
	hub.register <- ps

	defer func() {
		hub.unregister <- ps
		dc.Close()
		ps.Close()
	}()

	for {
		// handle PeerConnection events and close SSE event.
		select {
		case candidate := <-ps.candidateFound:
			jsonValue, err := messages.NewCandidateMessageJSON(newClient.userName, candidate)
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
			case webrtc.PeerConnectionStateConnected:
				for _, rcv := range peerConnection.GetReceivers() {
					track := rcv.Track()
					if track == nil {
						continue
					}
					log.Printf("RECV ID: %s MID: %s MSID: %s Kind: %s", track.ID(), track.RID(), track.Msid(), track.Kind())
				}
			case webrtc.PeerConnectionStateFailed:
				return
			case webrtc.PeerConnectionStateClosed:
				return
			}
		case message := <-ps.channels.MessageCh:
			if message.Error != nil {
				log.Printf("DataChannelError: %s\n", message.Error.Error())
				return
			}
			// Send back heartbeat message.
			if message.ID == HeartBeatChID {
				for id, dc := range ps.channels.DataChannels {
					if id == message.ID {
						dc.Send(message.Message.Data)
					}
				}
			} else {
				hub.broadcastDataChannelMessage <- receivedDataChannelMessage{
					UserName: userName,
					Message:  message,
				}
			}
		case <-r.Context().Done():
			// when "es.close()" is called, this loop operation will be ended.
			return
		}
	}
}
func rejectConnection(w http.ResponseWriter, r *http.Request, message string) {
	flusher, _ := w.(http.Flusher)
	_, _ = fmt.Fprintf(w, "data: %s\n\n", messages.NewErrorMessageJSON(message))
	flusher.Flush()
	for range r.Context().Done() {
		// when "es.close()" is called, this loop operation will be ended.
		return
	}
}
func SendSSEMessage(w http.ResponseWriter, hub *SSEHub, message *messages.ClientMessage) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	hub.broadcast <- *message
	data, _ := json.Marshal(GetSucceeded())
	w.Write(data)
}
