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
	printInfos(nodeInfos)
	fmt.Println("which do you want to link")
	targetId := 0
	fmt.Scanf("%d", &targetId)
	targetNode := findInfos(nodeInfos, targetId)
	if targetNode.Id == 0 {
		log.Panicln("not found node")
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
func findInfos(nodeInfos []NodeInfo, id int) NodeInfo {
	for _, v := range nodeInfos {
		if v.Id == id {
			return v
		}
	}
	return NodeInfo{Id: 0}
}
func printInfos(nodeInfos []NodeInfo) {
	for _, v := range nodeInfos {
		fmt.Printf("id: %03d | addr: %s | hostname: %s\n", v.Id, v.Addr, v.Hostname)
	}
}
