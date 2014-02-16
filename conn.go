package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	WRITE_WAIT  = 10 * time.Second
	PING_PERIOD = 45 * time.Second
)

type Conn struct {
	socket *websocket.Conn
	output chan []byte
}

func CreateConn(w http.ResponseWriter, r *http.Request) *Conn {
	socket, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return nil
	} else if err != nil {
		log.Println(err)
		return nil
	}
	return &Conn{output: make(chan []byte, 256), socket: socket}
}

// write writes a message with the given message type and payload.
func (c *Conn) write(messageType int, payload []byte) error {
	c.socket.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
	return c.socket.WriteMessage(messageType, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Conn) WritePump() {
	ticker := time.NewTicker(PING_PERIOD)
	defer func() {
		ticker.Stop()
		c.socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.output:
			if !ok {
				log.Println("Sending close message")
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			log.Println("Sending text")
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			log.Println("Sending ping")
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
