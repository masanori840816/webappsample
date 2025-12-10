package main

import (
	"log"
	"sync"

	"github.com/pion/webrtc/v4"
)

type PeerConnectionState struct {
	peerConnection        *webrtc.PeerConnection
	client                *SSEClient
	channels              *WebRTCDataChannelStates
	candidateFound        chan *webrtc.ICECandidate
	changeConnectionState chan webrtc.PeerConnectionState
	addTrack              chan *webrtc.TrackRemote
	heartbeat             chan int
	closeLock             sync.Mutex
	closed                bool
}

func (ps *PeerConnectionState) Close() {
	ps.closeLock.Lock()
	defer ps.closeLock.Unlock()
	if ps.closed {
		return
	}
	ps.closed = true
	for i := 0; i < len(ps.heartbeat); i++ {
		<-ps.heartbeat
	}
	close(ps.heartbeat)
	close(ps.candidateFound)
	close(ps.changeConnectionState)
	close(ps.addTrack)
	if ps.peerConnection.ConnectionState() != webrtc.PeerConnectionStateClosed {
		ps.peerConnection.Close()
	}
}
func NewPeerConnection() (*webrtc.PeerConnection, error) {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					"stun:stun.l.google.com:19302",
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			return nil, err
		}
	}
	return peerConnection, nil
}
func NewPeerConnectionState(client *SSEClient, peerConnection *webrtc.PeerConnection,
	channels *WebRTCDataChannelStates) (*PeerConnectionState, error) {
	heartbeat := make(chan int, 1)
	candidateFound := make(chan *webrtc.ICECandidate)
	changeConnectionState := make(chan webrtc.PeerConnectionState)
	addTrack := make(chan *webrtc.TrackRemote)

	heartbeat <- 1

	ps := &PeerConnectionState{
		peerConnection:        peerConnection,
		client:                client,
		channels:              channels,
		candidateFound:        candidateFound,
		changeConnectionState: changeConnectionState,
		addTrack:              addTrack,
		heartbeat:             heartbeat,
	}

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}
		ps.closeLock.Lock()
		defer ps.closeLock.Unlock()
		if ps.closed {
			return
		}
		_, ok := <-heartbeat
		if ok {
			candidateFound <- i
			heartbeat <- 1
		}
	})
	peerConnection.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		ps.closeLock.Lock()
		defer ps.closeLock.Unlock()
		if ps.closed {
			return
		}
		_, ok := <-heartbeat
		if ok {
			changeConnectionState <- p
			heartbeat <- 1
		}
		log.Printf("State: %s", p.String())
	})
	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		ps.closeLock.Lock()
		defer ps.closeLock.Unlock()
		if ps.closed {
			return
		}
		log.Printf("OnTrack ID:%s Kind:%s MSID:%s Codec:%s SSRC:%d StreamID:%s", t.ID(), t.Kind(), t.Msid(), t.Codec().MimeType, t.SSRC(), t.StreamID())
		_, ok := <-heartbeat
		if ok {
			addTrack <- t
			heartbeat <- 1
		}
	})

	return ps, nil
}
