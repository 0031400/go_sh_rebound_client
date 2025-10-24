package main

import (
	"encoding/json"
	"fmt"
	"go_sh_rebound_client/config"
	"go_sh_rebound_client/logger"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"golang.org/x/term"
)

type NodeInfo struct {
	Id       int    `json:"id"`
	Hostname string `json:"hostname"`
	Addr     string `json:"addr"`
}

func main() {
	logger.Init()
	config.Init()
	req, err := http.NewRequest("GET", config.ServerNodes, nil)
	if err != nil {
		log.Panicln(err)
	}
	req.Header.Set("authorization", config.Auth)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	d, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panicln(err)
	}
	var nodeInfos []NodeInfo
	err = json.Unmarshal(d, &nodeInfos)
	if err != nil {
		log.Fatalln(err)
	}
	targetId := chooseTheNodeId(nodeInfos)
	if targetId == 0 {
		log.Panicln("choose the node id fail")
	}
	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s?id=%d", config.ServerWs, targetId), http.Header{"authorization": []string{config.Auth}})
	if err != nil {
		log.Panicln(err)
	}
	defer c.Close()
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Panicln(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				log.Panicln(err)
			}
			if n > 0 {
				c.WriteMessage(websocket.BinaryMessage, buf[:n])
			}
		}
	}()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Panicln(err)
		}
		if mt == websocket.BinaryMessage {
			os.Stdout.Write(message)
		}
	}
}
func chooseTheNodeId(nodeInfos []NodeInfo) int {
	if len(nodeInfos) == 0 {
		fmt.Println("there is no node connecting with the server")
		return 0
	}
	for _, v := range nodeInfos {
		fmt.Printf("id: %03d | addr: %s | hostname: %s\n", v.Id, v.Addr, v.Hostname)
	}
	fmt.Println("which do you want to connect")
	targetId := 0
	fmt.Scanf("%d", &targetId)
	for _, v := range nodeInfos {
		if v.Id == targetId {
			return targetId
		}
	}
	return 0
}
