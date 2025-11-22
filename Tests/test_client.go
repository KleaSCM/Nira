/**
 * Simple test client for NIRA backend.
 *
 * Connects to the WebSocket server and sends a test message to verify
 * the Ollama integration is working correctly.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 * File: test_client.go
 * Description: Test client for backend verification.
 */

package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Dial error: %v", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("Read error: %v", err)
				return
			}
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("JSON parse error: %v", err)
				continue
			}
			log.Printf("Received: type=%v, content=%v", msg["type"], msg["content"])
		}
	}()

	testMsg := map[string]interface{}{
		"type":    "user",
		"content": "Hello! Can you tell me a short joke?",
	}

	msgBytes, err := json.Marshal(testMsg)
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}

	log.Println("Sending test message...")
	err = c.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		log.Fatalf("Write error: %v", err)
	}

	timeout := time.After(30 * time.Second)
	for {
		select {
		case <-done:
			return
		case <-timeout:
			log.Println("Test timeout")
			return
		case <-interrupt:
			log.Println("Interrupt received")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("Write close error: %v", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
