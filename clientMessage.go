package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/pion/webrtc/v3"
)

const (
	TextEvent       string = "text"
	OfferEvent      string = "offer"
	AnswerEvent     string = "answer"
	CandidateEvent  string = "candidate"
	UpdateEvent     string = "update"
	ClientNameEvent string = "clientName"
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
	message, err := NewOfferMessage(userName, offer)
	if err != nil {
		return "", err
	}
	jsonValue, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(jsonValue), nil
}
func NewCandidateMessage(userName string, candidate *webrtc.ICECandidate) (*ClientMessage, error) {
	if candidate == nil {
		return nil, errors.New("ICECandidate was null")
	}

	candidateString, err := json.Marshal(candidate.ToJSON())
	if err != nil {
		return nil, err
	}
	return &ClientMessage{
		Event:    CandidateEvent,
		UserName: userName,
		Data:     string(candidateString),
	}, nil
}
func NewCandidateMessageJSON(userName string, candidate *webrtc.ICECandidate) (string, error) {
	message, err := NewCandidateMessage(userName, candidate)
	if err != nil {
		return "", err
	}
	jsonValue, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	log.Printf("NewCandidate: %s", string(jsonValue))
	return string(jsonValue), nil
}
func NewClientNameMessageJSON(names ClientNames) (string, error) {
	clientNamesJson, err := json.Marshal(names)
	if err != nil {
		return "", err
	}
	message := ClientMessage{
		Event:    ClientNameEvent,
		UserName: "",
		Data:     string(clientNamesJson),
	}
	jsonValue, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(jsonValue), nil
}
