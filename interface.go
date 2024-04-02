package main

import "github.com/google/uuid"

var globalStatus = GlobalStatus{}
var Clients = make(map[uuid.UUID]*Client)

type GlobalStatus struct {
	Second        int
	InputSecond   int
	IsPause       bool
	IsStart       bool
	StopwatchMode bool
}

type Client struct {
	OutCh chan string
}
