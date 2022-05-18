package main

type ClientMessage struct {
	Event    string `json:"event"`
	UserName string `json:"userName"`
	Data     string `json:"data"`
}
