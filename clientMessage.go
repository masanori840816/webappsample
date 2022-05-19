package main

import (
	"encoding/json"

	"github.com/pion/webrtc/v3"
)

const (
	TextEvent      string = "text"
	OfferEvent     string = "offer"
	AnswerEvent    string = "answer"
	CandidateEvent string = "candidate"
)

type ClientMessage struct {
	Event    string `json:"event"`
	UserName string `json:"userName"`
	Data     string `json:"data"`
}

func NewOfferMessage(userName string, offer webrtc.SessionDescription) (*ClientMessage, error) {
	offerString, err := json.Marshal(offer)
	if err != nil {
		return nil, err
	}
	return &ClientMessage{
		Event:    OfferEvent,
		UserName: userName,
		Data:     string(offerString),
	}, nil
}
func NewOfferMessageJSON(userName string, offer webrtc.SessionDescription) (string, error) {
	offerMessage, err := NewOfferMessage(userName, offer)
	if err != nil {
		return "", err
	}
	offerJSON, err := json.Marshal(offerMessage)
	if err != nil {
		return "", err
	}
	return string(offerJSON), nil
}
