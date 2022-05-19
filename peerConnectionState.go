package main

import (
	"log"

	"github.com/pion/webrtc/v3"
)

type PeerConnectionState struct {
	peerConnection *webrtc.PeerConnection
	client         *SSEClient
}

func NewPeerConnectionState(client *SSEClient) (*PeerConnectionState, error) {
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
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}
		log.Println("OnICECandidate")
		_, ok := <-client.candidateChan
		if ok {
			client.candidateChan <- i
		}
	})
	peerConnection.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		log.Println("OnConnectionStateChange")
		_, ok := <-client.changeConnectionState
		if ok {
			client.changeConnectionState <- p
		}
	})
	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		log.Println("OnTrack")
		_, ok := <-client.addTrack
		if ok {
			client.addTrack <- t
		}
	})

	return &PeerConnectionState{
		peerConnection: peerConnection,
		client:         client,
	}, nil
}
