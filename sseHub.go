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
	addTrack    chan *webrtc.TrackRemote
}

func newSSEHub() *SSEHub {
	return &SSEHub{
		clients:     make(map[*PeerConnectionState]bool),
		broadcast:   make(chan ClientMessage),
		register:    make(chan *PeerConnectionState),
		unregister:  make(chan *PeerConnectionState),
		trackLocals: map[string]*webrtc.TrackLocalStaticRTP{},
		addTrack:    make(chan *webrtc.TrackRemote),
	}
}
func newTrackLocalStaticRTP(track *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error) {
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
					log.Println("close after unregister")
					client.peerConnection.Close()
				}
				delete(h.clients, client)
				signalPeerConnections(h)
			}
		case track := <-h.addTrack:
			trackLocal, err := newTrackLocalStaticRTP(track)
			if err != nil {
				log.Println(err.Error())
				return
			}
			h.trackLocals[track.ID()] = trackLocal
			signalPeerConnections(h)
			defer func() {
				log.Println("delete tracks")
				delete(h.trackLocals, track.ID())
				signalPeerConnections(h)
			}()
			go updateTrackValue(track, trackLocal)

		case message := <-h.broadcast:
			log.Println("MMMMEssage")
			handleReceivedMessage(h, message)
		}
	}
}
func updateTrackValue(track *webrtc.TrackRemote, trackLocal *webrtc.TrackLocalStaticRTP) {
	buf := make([]byte, 1500)

	for {
		i, _, err := track.Read(buf)
		if err != nil {
			return
		}

		if _, err = trackLocal.Write(buf[:i]); err != nil {
			return
		}
	}
}
func handleReceivedMessage(h *SSEHub, message ClientMessage) {
	switch message.Event {
	case TextEvent:
		m, _ := json.Marshal(message)
		jsonText := string(m)

		for client := range h.clients {
			select {
			case client.client.sendMessage <- jsonText:
			default:
				client.client.CloseAllChannels()
				if client.peerConnection.ConnectionState() == webrtc.PeerConnectionStateConnected {
					client.peerConnection.Close()
				}
				delete(h.clients, client)
			}
		}
	case CandidateEvent:
		candidate := webrtc.ICECandidateInit{}
		if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
			log.Println(err)
			return
		}
		for pc := range h.clients {
			if pc.client.userName == message.UserName {
				if err := pc.peerConnection.AddICECandidate(candidate); err != nil {
					log.Println(err)
					return
				}
			}
		}
	case AnswerEvent:
		answer := webrtc.SessionDescription{}
		if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
			log.Println(err)
			return
		}
		for pc := range h.clients {
			if pc.client.userName == message.UserName {
				if err := pc.peerConnection.SetRemoteDescription(answer); err != nil {
					log.Println(err)
					return
				}
			}
		}

	}
}
func signalPeerConnections(h *SSEHub) {
	log.Println("signalPeerConnections")
	for ps := range h.clients {
		if ps.peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
			delete(h.clients, ps)
			// We modified the slice, start from the beginning
			signalPeerConnections(h)

			log.Println("state closed")
		}
		log.Println("signalPeerConnections1")

		existingSenders := map[string]bool{}

		for _, sender := range ps.peerConnection.GetSenders() {
			if sender.Track() == nil {
				continue
			}
			log.Println("sender " + sender.Track().ID())
			existingSenders[sender.Track().ID()] = true

			if _, ok := h.trackLocals[sender.Track().ID()]; !ok {
				if err := ps.peerConnection.RemoveTrack(sender); err != nil {
					log.Println(err.Error())
					return
				}
			}
		}
		log.Println("signalPeerConnections3")

		for _, receiver := range ps.peerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}

			existingSenders[receiver.Track().ID()] = true
		}
		log.Println("signalPeerConnections4")

		for trackID := range h.trackLocals {
			if _, ok := existingSenders[trackID]; !ok {
				if _, err := ps.peerConnection.AddTrack(h.trackLocals[trackID]); err != nil {
					log.Println(err.Error())
					return
				}
			}
		}
		log.Println("signalPeerConnections555")

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

		log.Println("signalPeerConnections3")
		ps.client.sendMessage <- messageJSON
		log.Println("signalPeerConnections4")
	}
}
