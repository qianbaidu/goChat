package main

import (
	"fmt"
	"github.com/prometheus/common/log"
	"golang.org/x/net/websocket"
)

const url = "ws://127.0.0.1:8080/chat?userName=hello"
const origin = "http://127.0.0.1:8080/"

func main() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	message := []byte("hello, world")
	_, err = ws.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", message)

	var msg = make([]byte, 512)
	m, err := ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg[:m])

	ws.Close()
}
