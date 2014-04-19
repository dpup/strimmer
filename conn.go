// Copyright 2014 Daniel Pupius
// Based on tutorial at http://gary.burd.info/go-websocket-chat

package main

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	WRITE_WAIT  = 10 * time.Second
	PING_PERIOD = 15 * time.Second
)

type conn struct {
	socket  *websocket.Conn
	output  chan []byte
	created time.Time
}

// write writes a message with the given message type and payload.
func (c *conn) write(messageType int, payload []byte) error {
	c.socket.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
	return c.socket.WriteMessage(messageType, payload)
}

// wait handles messages sent from the hub on the connection's channel to
// the websocket connection.
func (c *conn) wait() {
	ticker := time.NewTicker(PING_PERIOD)
	defer func() {
		ticker.Stop()
		c.socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.output:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
