package main

import (
	"encoding/json"

	"github.com/pion/webrtc/v3"
)

type SSEHub struct {
	clients     map[*PeerConnectionState]bool
	broadcast   chan ClientMessage
	register    chan *PeerConnectionState
	unregister  chan *PeerConnectionState
	trackLocals map[string]*webrtc.TrackLocalStaticRTP
	addTrack    chan *webrtc.TrackLocalStaticRTP
	removeTrack chan *webrtc.TrackLocalStaticRTP
}

func newSSEHub() *SSEHub {
	return &SSEHub{
		clients:     make(map[*PeerConnectionState]bool),
		broadcast:   make(chan ClientMessage),
		register:    make(chan *PeerConnectionState),
		unregister:  make(chan *PeerConnectionState),
		trackLocals: map[string]*webrtc.TrackLocalStaticRTP{},
		addTrack:    make(chan *webrtc.TrackLocalStaticRTP),
		removeTrack: make(chan *webrtc.TrackLocalStaticRTP),
	}
}
func generateTrackLocalStaticRTP(track *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error) {
	return webrtc.NewTrackLocalStaticRTP(track.Codec().RTPCodecCapability, track.ID(), track.StreamID())
}
func (h *SSEHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				client.client.CloseAllChannels()
				if client.peerConnection.ConnectionState() == webrtc.PeerConnectionStateConnected {
					client.peerConnection.Close()
				}
				delete(h.clients, client)
			}
		case track := <-h.addTrack:
			h.trackLocals[track.ID()] = track
		case track := <-h.removeTrack:
			delete(h.trackLocals, track.ID())
		case message := <-h.broadcast:
			m, _ := json.Marshal(message)
			jsonText := string(m)

			for client := range h.clients {
				if client.client.userName == message.UserName {
					continue
				}
				select {
				case client.client.messageChan <- jsonText:
				default:
					client.client.CloseAllChannels()
					if client.peerConnection.ConnectionState() == webrtc.PeerConnectionStateConnected {
						client.peerConnection.Close()
					}
					delete(h.clients, client)
				}
			}
		}
	}
}
