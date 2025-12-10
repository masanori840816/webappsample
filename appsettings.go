package main

import (
	"encoding/json"
	"os"
)

type AppSettings struct {
	URL           string `json:"url"`
	ICEServerJSON string `json:"iceServer"`
}
type ICEServer struct {
	URLs       string `json:"urls"`
	UserName   string `json:"username"`
	Credential string `json:"credential"`
}

func LoadAppSettings() (setting AppSettings, err error) {
	result := &AppSettings{}
	result.URL = os.Getenv("WEBRTCAPP_URL")
	iceServer := &ICEServer{}
	iceServer.URLs = os.Getenv("WEBRTCAPP_ICE_URL")
	iceServer.UserName = os.Getenv("WEBRTCAPP_ICE_USERNAME")
	iceServer.Credential = os.Getenv("WEBRTCAPP_ICE_CREDENTIAL")
	if len(result.URL) <= 0 {
		result.URL = "http://localhost:8080"
	}
	if len(iceServer.URLs) <= 0 {
		iceServer.URLs = "stun:stun.l.google.com:19302"
	}
	iceJSON, err := json.Marshal(iceServer)
	if err != nil {
		return AppSettings{}, err
	}
	result.ICEServerJSON = string(iceJSON)
	return *result, err
}
