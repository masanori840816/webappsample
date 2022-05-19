package main

import (
	"encoding/json"
	"log"

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
			signalPeerConnections(h)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				client.client.CloseAllChannels()
				if client.peerConnection.ConnectionState() == webrtc.PeerConnectionStateConnected {
					client.peerConnection.Close()
				}
				delete(h.clients, client)
				signalPeerConnections(h)
			}
		case track := <-h.addTrack:
			h.trackLocals[track.ID()] = track
			signalPeerConnections(h)
		case track := <-h.removeTrack:
			delete(h.trackLocals, track.ID())
			signalPeerConnections(h)
		case message := <-h.broadcast:
			m, _ := json.Marshal(message)
			jsonText := string(m)

			for client := range h.clients {
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
func signalPeerConnections(h *SSEHub) {
	for ps := range h.clients {
		if ps.peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
			delete(h.clients, ps)
			// We modified the slice, start from the beginning
			signalPeerConnections(h)
		}

		existingSenders := map[string]bool{}

		for _, sender := range ps.peerConnection.GetSenders() {
			if sender.Track() == nil {
				continue
			}

			existingSenders[sender.Track().ID()] = true

			if _, ok := h.trackLocals[sender.Track().ID()]; !ok {
				if err := ps.peerConnection.RemoveTrack(sender); err != nil {
					log.Println(err.Error())
					return
				}
			}
		}

		for _, receiver := range ps.peerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}

			existingSenders[receiver.Track().ID()] = true
		}

		for trackID := range h.trackLocals {
			if _, ok := existingSenders[trackID]; !ok {
				if _, err := ps.peerConnection.AddTrack(h.trackLocals[trackID]); err != nil {
					log.Println(err.Error())
					return
				}
			}
		}

		offer, err := ps.peerConnection.CreateOffer(nil)
		if err != nil {
			log.Println(err.Error())
			return
		}

		if err = ps.peerConnection.SetLocalDescription(offer); err != nil {
			log.Println(err.Error())
			return
		}
		messageJSON, err := NewOfferMessageJSON(ps.client.userName, offer)
		if err != nil {
			log.Println(err.Error())
			return
		}
		ps.client.messageChan <- messageJSON
	}
}
