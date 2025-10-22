package config

import (
	"flag"
)

var ServerWs = "ws://127.0.0.1:3000/client/ws"
var ServerNodes = "http://127.0.0.1:3000/nodes"
var Auth = ""

func Init() {
	wsStr := flag.String("ws", "", "the server ws addr")
	nodesStr := flag.String("nodes", "", "the server nodes addr")
	authStr := flag.String("a", "", "the authorization header")
	flag.Parse()
	if *wsStr != "" {
		ServerWs = *wsStr
	}
	if *nodesStr != "" {
		ServerNodes = *nodesStr
	}
	if *authStr != "" {
		Auth = *authStr
	}
}
