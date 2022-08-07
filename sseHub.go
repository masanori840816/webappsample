package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
)

type SSEHub struct {
	clients                     map[*PeerConnectionState]bool
	broadcast                   chan ClientMessage
	register                    chan *PeerConnectionState
	unregister                  chan *PeerConnectionState
	trackLocals                 map[string]*webrtc.TrackLocalStaticRTP
	addTrack                    chan *webrtc.TrackRemote
	broadcastDataChannelMessage chan ReceivedDataChannelMessage
}

func newSSEHub() *SSEHub {
	return &SSEHub{
		clients:                     make(map[*PeerConnectionState]bool),
		broadcast:                   make(chan ClientMessage),
		register:                    make(chan *PeerConnectionState),
		unregister:                  make(chan *PeerConnectionState),
		trackLocals:                 map[string]*webrtc.TrackLocalStaticRTP{},
		addTrack:                    make(chan *webrtc.TrackRemote),
		broadcastDataChannelMessage: make(chan ReceivedDataChannelMessage),
	}
}
func (h *SSEHub) close() {
	close(h.broadcast)
	close(h.register)
	close(h.unregister)
	close(h.addTrack)
	close(h.broadcastDataChannelMessage)
}
func (h *SSEHub) run() {
	defer h.close()
	go func() {
		for range time.NewTicker(time.Second * 3).C {
			dispatchKeyFrame(h)
		}
	}()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			signalPeerConnections(h)
			sendClientNames(h)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				signalPeerConnections(h)
				sendClientNames(h)
			}
		case track := <-h.addTrack:
			trackLocal, err := webrtc.NewTrackLocalStaticRTP(track.Codec().RTPCodecCapability,
				track.ID(), track.StreamID())
			if err != nil {
				log.Println(err.Error())
				return
			}
			h.trackLocals[track.ID()] = trackLocal
			signalPeerConnections(h)
			go updateTrackValue(h, track)

		case message := <-h.broadcast:
			handleReceivedMessage(h, message)
		case message := <-h.broadcastDataChannelMessage:
			sendDataChannelMessage(h, message)
		}
	}
}
func updateTrackValue(h *SSEHub, track *webrtc.TrackRemote) {
	log.Printf("updateTrackValue Track: %s Kind: %s MSID: %s", track.ID(), track.Kind(), track.Msid())
	defer func() {
		delete(h.trackLocals, track.ID())
		signalPeerConnections(h)
	}()

	buf := make([]byte, 1500)

	for {
		i, _, err := track.Read(buf)
		if err != nil {
			return
		}
		if _, err = h.trackLocals[track.ID()].Write(buf[:i]); err != nil {
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
			flusher, _ := client.client.w.(http.Flusher)

			fmt.Fprintf(client.client.w, "data: %s\n\n", jsonText)
			flusher.Flush()
		}
	case CandidateEvent:
		log.Printf("R Candidate: %s", message.Data)
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
		log.Printf("R Answer: %s", message.Data)
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
	case UpdateEvent:
		signalPeerConnections(h)
	}
}
func signalPeerConnections(h *SSEHub) {
	defer func() {
		dispatchKeyFrame(h)
	}()
	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			// Release the lock and attempt a sync in 3 seconds. We might be blocking a RemoveTrack or AddTrack
			go func() {
				time.Sleep(time.Second * 3)
				signalPeerConnections(h)
			}()
			return
		}

		if !attemptSync(h) {
			break
		}
	}
}
func attemptSync(h *SSEHub) bool {
	for ps := range h.clients {
		if ps.peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
			delete(h.clients, ps)
			// We modified the slice, start from the beginning
			return true
		}
		existingSenders := map[string]bool{}

		for _, sender := range ps.peerConnection.GetSenders() {
			if sender.Track() == nil {
				continue
			}
			existingSenders[sender.Track().ID()] = true

			if _, ok := h.trackLocals[sender.Track().ID()]; !ok {
				if err := ps.peerConnection.RemoveTrack(sender); err != nil {
					return true
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
					return true
				}
			}
		}

		offer, err := ps.peerConnection.CreateOffer(nil)
		if err != nil {
			return true
		}
		messageJSON, err := NewOfferMessageJSON(ps.client.userName, offer)
		if err != nil {
			return true
		}

		if err = ps.peerConnection.SetLocalDescription(offer); err != nil {
			return true
		}
		flusher, _ := ps.client.w.(http.Flusher)
		log.Println(messageJSON)
		fmt.Fprintf(ps.client.w, "data: %s\n\n", messageJSON)
		flusher.Flush()
	}
	return false
}
func dispatchKeyFrame(h *SSEHub) {
	for ps := range h.clients {
		for _, receiver := range ps.peerConnection.GetReceivers() {
			if receiver.Track() == nil {
				log.Printf("dispatchKeyFrame Track nil C: %s", ps.client.userName)
				continue
			}
			t := receiver.Track()
			log.Printf("dispatchKeyFrame C: %s TID: %s Kind: %s MSID: %s", ps.client.userName, t.ID(), t.Kind(), t.Msid())

			_ = ps.peerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}
func sendDataChannelMessage(h *SSEHub, message ReceivedDataChannelMessage) {
	for ps := range h.clients {
		if ps.client.userName != message.UserName {
			ps.channels.SendMessage(message)
		}
	}
}
func sendClientNames(h *SSEHub) {
	names := ClientNames{
		Names: make([]ClientName, len(h.clients)),
	}

	i := 0
	for ps := range h.clients {
		names.Names[i] = ClientName{
			Name: ps.client.userName,
		}
	}
	message, err := NewClientNameMessageJSON(names)
	if err != nil {
		log.Printf("Error sendClientNames Message: %s", err.Error())
		return
	}
	for ps := range h.clients {
		flusher, _ := ps.client.w.(http.Flusher)
		fmt.Fprintf(ps.client.w, "data: %s\n\n", message)
		flusher.Flush()
	}
}
