package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/pion/webrtc/v4"
)

const (
	TextEvent       string = "text"
	OfferEvent      string = "offer"
	AnswerEvent     string = "answer"
	CandidateEvent  string = "candidate"
	UpdateEvent     string = "update"
	ClientNameEvent string = "clientName"
	HeartbeatEvent  string = "heartbeat"
	ErrorEvent      string = "error"
)

type ClientMessage struct {
	MessageType string `json:"event"`
	UserName    string `json:"userName"`
	GroupName   string `json:"groupName"`
	Data        string `json:"data"`
}

func NewOfferMessage(userName string, offer webrtc.SessionDescription) (*ClientMessage, error) {
	offerString, err := json.Marshal(offer)
	if err != nil {
		return nil, err
	}
	return &ClientMessage{
		MessageType: OfferEvent,
		UserName:    userName,
		Data:        string(offerString),
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
		MessageType: CandidateEvent,
		UserName:    userName,
		Data:        string(candidateString),
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
		MessageType: ClientNameEvent,
		UserName:    "",
		Data:        string(clientNamesJson),
	}
	jsonValue, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(jsonValue), nil
}

func NewHeartbeatMessageJSON() string {
	message := ClientMessage{
		MessageType: HeartbeatEvent,
		Data:        ":",
	}
	jsonValue, err := json.Marshal(message)
	if err != nil {
		return ""
	}
	return string(jsonValue)
}

func NewErrorMessageJSON(message string) string {
	faildMessage := ClientMessage{
		MessageType: ErrorEvent,
		Data:        message,
	}
	resultJson, err := json.Marshal(faildMessage)
	if err != nil {
		return ""
	}
	return string(resultJson)
}
