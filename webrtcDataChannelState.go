package main

import (
	"log"

	"github.com/pion/webrtc/v3"
)

type WebRTCDataChannelStates struct {
	DataChannels map[uint16]*webrtc.DataChannel
	MessageCh    chan WebRTCDataChannelMessage
}
type WebRTCDataChannelMessage struct {
	ID      uint16
	Message webrtc.DataChannelMessage
	Error   error
}
type ReceivedDataChannelMessage struct {
	UserName string
	ID       uint16
	Message  webrtc.DataChannelMessage
}

func (dc *WebRTCDataChannelStates) SendMessage(message ReceivedDataChannelMessage) {
	for i, channel := range dc.DataChannels {
		if i == message.ID {
			channel.Send(message.Message.Data)
		}
	}
}
func (dc *WebRTCDataChannelStates) Close() {
	for _, v := range dc.DataChannels {
		v.Close()
	}
	close(dc.MessageCh)
}
func NewWebRTCDataChannelStates(peerConnection *webrtc.PeerConnection) (*WebRTCDataChannelStates, error) {
	dc := &WebRTCDataChannelStates{
		DataChannels: make(map[uint16]*webrtc.DataChannel),
		MessageCh:    make(chan WebRTCDataChannelMessage),
	}
	err := addWebRTCDataChannel("sample", 20, peerConnection, dc)
	if err != nil {
		return nil, err
	}
	err = addWebRTCDataChannel("sample2", 21, peerConnection, dc)
	if err != nil {
		return nil, err
	}
	return dc, nil
}
func addWebRTCDataChannel(label string, id uint16, peerConnection *webrtc.PeerConnection,
	dc *WebRTCDataChannelStates) error {
	dataChannel, err := newWebRTCDataChannel(label, id, peerConnection)
	if err != nil {
		return err
	}
	dc.DataChannels[*dataChannel.ID()] = dataChannel
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		dc.MessageCh <- WebRTCDataChannelMessage{
			ID:      id,
			Message: msg,
		}
	})
	dataChannel.OnError(func(err error) {
		dc.MessageCh <- WebRTCDataChannelMessage{
			ID:    id,
			Error: err,
		}
	})
	return nil
}
func newWebRTCDataChannel(label string, id uint16, peerConnection *webrtc.PeerConnection) (*webrtc.DataChannel, error) {
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

	dc.OnClose(func() {
		log.Println("OnClose")
	})
	return dc, nil
}
