package main

import (
	"os"
)

type AppSettings struct {
	URL string `json:"url"`
}

func LoadAppSettings() (setting AppSettings, err error) {
	result := &AppSettings{}
	result.URL = os.Getenv("WEBRTCAPP_URL")
	if len(result.URL) <= 0 {
		result.URL = "http://localhost:8080"
	}
	return *result, err
}
