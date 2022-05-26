package main

import (
	"log"

	"github.com/pion/webrtc/v3"
)

type WebRTCDataChannelState struct {
	DataChannel *webrtc.DataChannel
}

func (dc *WebRTCDataChannelState) Close() {
	dc.DataChannel.Close()
}
func NewWebRTCDataChannelState(label string, id uint16, peerConnection *webrtc.PeerConnection) (*WebRTCDataChannelState, error) {
	negotiated := true
	ordered := true
	dc, err := peerConnection.CreateDataChannel(label, &webrtc.DataChannelInit{
		ID:         &id,
		Negotiated: &negotiated,
		Ordered:    &ordered,
	})
	if err != nil {
		return nil, err
	}
	dc.OnOpen(func() {
		log.Println("OnOpen")
	})
	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		log.Println("OnMessage")
		log.Println(msg)
	})
	dc.OnClose(func() {
		log.Println("OnClose")
	})
	return &WebRTCDataChannelState{
		DataChannel: dc,
	}, nil
}
