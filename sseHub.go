package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"

	messages "webappsample/messages"
)

type SSEHub struct {
	name                        string
	clients                     map[*PeerConnectionState]bool
	broadcast                   chan messages.ClientMessage
	register                    chan *PeerConnectionState
	unregister                  chan *PeerConnectionState
	trackLocals                 map[string]*webrtc.TrackLocalStaticRTP
	addTrack                    chan *webrtc.TrackRemote
	broadcastDataChannelMessage chan receivedDataChannelMessage
}
type receivedDataChannelMessage struct {
	UserName string
	Message  WebRTCDataChannelMessage
}

func NewSSEHub(name string) *SSEHub {
	return &SSEHub{
		name:                        name,
		clients:                     make(map[*PeerConnectionState]bool),
		broadcast:                   make(chan messages.ClientMessage),
		register:                    make(chan *PeerConnectionState),
		unregister:                  make(chan *PeerConnectionState),
		trackLocals:                 map[string]*webrtc.TrackLocalStaticRTP{},
		addTrack:                    make(chan *webrtc.TrackRemote),
		broadcastDataChannelMessage: make(chan receivedDataChannelMessage),
	}
}
func (h *SSEHub) CloseSSEHub() {
	close(h.broadcast)
	close(h.register)
	close(h.unregister)
	close(h.addTrack)
	close(h.broadcastDataChannelMessage)
}
func (h *SSEHub) run(unregister chan *SSEHub) {
	heartbeatMessage := messages.NewHeartbeatMessageJSON()
	keyFrameTicker := time.NewTicker(time.Second * 3)
	heartbeat := time.NewTicker(time.Minute)
	defer func() {
		keyFrameTicker.Stop()
		heartbeat.Stop()
		unregister <- h
	}()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			signalPeerConnections(h)
			sendCurrentClientNames(h)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				signalPeerConnections(h)
				// Delete the group when there are no more clients
				if len(h.clients) <= 0 {
					return
				}
			}
			// send connected client names
			sendCurrentClientNames(h)

		case track := <-h.addTrack:
			trackLocal, err := webrtc.NewTrackLocalStaticRTP(track.Codec().RTPCodecCapability,
				track.ID(), track.StreamID())
			if err != nil {
				log.Printf("AddTrackError: %s", err.Error())
				return
			}
			h.trackLocals[track.ID()] = trackLocal
			signalPeerConnections(h)
			go updateTrackValue(h, track)

		case message := <-h.broadcast:
			handleReceivedMessage(h, message)
		case message := <-h.broadcastDataChannelMessage:
			for pc := range h.clients {
				if pc.client.userName != message.UserName {
					for id, dc := range pc.channels.DataChannels {
						if id == message.Message.ID {
							dc.Send(message.Message.Message.Data)
						}
					}
				}
			}
		case <-keyFrameTicker.C:
			dispatchKeyFrame(h)
		case <-heartbeat.C:
			for pc := range h.clients {
				flusher, _ := pc.client.w.(http.Flusher)
				fmt.Fprintf(pc.client.w, "data: %s\n\n", heartbeatMessage)
				flusher.Flush()
			}
		}
	}
}
func updateTrackValue(h *SSEHub, track *webrtc.TrackRemote) {
	log.Printf("updateTrackValue Track: %s Kind: %s MSID: %s MIME TYPE: %s", track.ID(), track.Kind(), track.Msid(), track.Codec().MimeType)
	defer func() {
		delete(h.trackLocals, track.ID())
		signalPeerConnections(h)
	}()

	buf := make([]byte, 1048576)

	for {
		i, _, err := track.Read(buf)
		if err != nil {
			log.Printf("updateTrackValue FailedReading Err: %s", err.Error())
			return
		}
		if _, ok := h.trackLocals[track.ID()]; !ok {
			log.Printf("updateTrackValue trackLocals doesn't have TID: %s", track.ID())
			return
		}
		if _, err = h.trackLocals[track.ID()].Write(buf[:i]); err != nil {
			log.Printf("updateTrackValue FailedWriting Err: %s", err.Error())
			return
		}
	}
}
func handleReceivedMessage(h *SSEHub, message messages.ClientMessage) {
	switch message.MessageType {
	case messages.TextEvent:
		m, _ := json.Marshal(message)
		jsonText := string(m)

		for client := range h.clients {
			flusher, _ := client.client.w.(http.Flusher)
			fmt.Fprintf(client.client.w, "data: %s\n\n", jsonText)
			flusher.Flush()
		}
	case messages.CandidateEvent:
		candidate, err := parseICECandidate(message.Data)
		if err != nil {
			log.Println(err)
			return
		}
		for pc := range h.clients {
			if pc.client.userName == message.UserName {
				if err := pc.peerConnection.AddICECandidate(candidate); err != nil {
					log.Printf("AddICECandidate erorr:%s", err.Error())
					return
				}
			}
		}
	case messages.AnswerEvent:
		answer := webrtc.SessionDescription{}
		if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
			log.Println(err)
			return
		}
		for pc := range h.clients {
			if pc.client.userName == message.UserName {
				if err := pc.peerConnection.SetRemoteDescription(answer); err != nil {
					log.Printf("SetRemoteDesc Error: %s", err.Error())
					return
				}
			}
		}
	}
}
func parseICECandidate(value string) (webrtc.ICECandidateInit, error) {
	splittedData := strings.Split(value, "|")
	if len(splittedData) < 3 {
		log.Println(value)
		return webrtc.ICECandidateInit{}, errors.New("Failed to read ICE Candidate")
	}
	lineIndex64, err := strconv.ParseUint(splittedData[1], 10, 16)
	if err != nil {
		return webrtc.ICECandidateInit{}, err
	}
	lineIndex := uint16(lineIndex64)
	candidate := webrtc.ICECandidateInit{
		SDPMLineIndex: &lineIndex,
		SDPMid:        &splittedData[2],
		Candidate:     splittedData[0],
	}
	return candidate, nil
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
		if ps.peerConnection.ConnectionState() != webrtc.PeerConnectionStateConnected &&
			ps.peerConnection.ConnectionState() != webrtc.PeerConnectionStateConnecting &&
			ps.peerConnection.ConnectionState() != webrtc.PeerConnectionStateNew {
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
		messageJSON, err := messages.NewOfferMessageJSON(ps.client.userName, offer)
		if err != nil {
			return true
		}

		if err = ps.peerConnection.SetLocalDescription(offer); err != nil {
			return true
		}
		flusher, _ := ps.client.w.(http.Flusher)
		log.Printf("Offer C: %s V: %s", ps.client.userName, messageJSON)
		fmt.Fprintf(ps.client.w, "data: %s\n\n", messageJSON)
		flusher.Flush()
	}
	return false
}
func dispatchKeyFrame(h *SSEHub) {
	for ps := range h.clients {
		for _, receiver := range ps.peerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}
			_ = ps.peerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}
func sendCurrentClientNames(h *SSEHub) {
	names := messages.ClientNames{
		Names: make([]messages.ClientName, len(h.clients)),
	}

	i := 0
	for pc := range h.clients {
		names.Names[i] = messages.ClientName{
			Name: pc.client.userName,
		}
		i += 1
	}
	message, err := messages.NewClientNameMessageJSON(names)
	if err != nil {
		log.Printf("Error sendClientNames Message: %s", err.Error())
		return
	}
	for pc := range h.clients {
		flusher, _ := pc.client.w.(http.Flusher)
		fmt.Fprintf(pc.client.w, "data: %s\n\n", message)
		flusher.Flush()
	}
}
