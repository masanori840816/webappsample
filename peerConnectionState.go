package main

import "github.com/pion/webrtc/v3"

type PeerConnectionState struct {
	peerConnection *webrtc.PeerConnection
	client         *SSEClient
}
