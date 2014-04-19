// Copyright 2014 Daniel Pupius
// Based on tutorial at http://gary.burd.info/go-websocket-chat

package main

import (
	"container/ring"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// TODO: Keep track of last message sent to each connection and send recent
// history on reconnect.

// Hub keeps a map of connections and broadcasts data to them.
type Hub struct {
	conns            map[*conn]int
	mu               sync.Mutex // protects conns
	connCounter      int
	history          *ring.Ring
	allowCrossOrigin bool
}

// Creates a new hub for managing websocket connections.
func NewHub(historySize int, allowCrossOrigin bool) *Hub {
	return &Hub{
		conns:            make(map[*conn]int),
		history:          ring.New(historySize),
		allowCrossOrigin: allowCrossOrigin,
	}
}

// HandleConnection responds to a HTTP request and attempts to upgrade the
// request into a websocket connection.
func (h *Hub) HandleConnection(w http.ResponseWriter, r *http.Request) {
	socket, err, statusCode := h.upgrade(w, r)
	if err != nil {
		log.Println("Error handling conncetion", err)
		http.Error(w, err.Error(), statusCode)
		return
	}

	conn := &conn{
		output:  make(chan []byte, 256),
		socket:  socket,
		created: time.Now(),
	}

	// Send recent history to clients when they connect.
	h.history.Do(func(v interface{}) {
		if v != nil {
			if data, ok := v.([]byte); ok {
				conn.output <- data
			}
		}
	})

	h.add(conn)

	defer func() { h.remove(conn) }()
	conn.wait()
}

func (h *Hub) upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error, int) {
	if r.Method != "GET" {
		return nil, errors.New("method not allowed"), http.StatusMethodNotAllowed
	}

	if !h.allowCrossOrigin && r.Header.Get("Origin") != "http://"+r.Host {
		return nil, errors.New("origin not allowed"), http.StatusMethodNotAllowed
	}

	socket, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		return nil, errors.New("not a websockethandshake"), http.StatusBadRequest
	}

	if err != nil {
		return nil, fmt.Errorf("error upgrading socket %s", err), http.StatusInternalServerError
	}

	return socket, nil, http.StatusOK
}

// Broadcast sends data to each active connection.
func (h *Hub) Broadcast(data []byte) {
	defer h.lock()()

	for c, id := range h.conns {
		select {
		case c.output <- data:
		default:
			log.Println("Connection closed :", id)
			delete(h.conns, c)
			close(c.output)
		}
	}

	// Store the history.
	h.history.Value = data
	h.history = h.history.Next()
}

// add a new  connection.
func (h *Hub) add(c *conn) {
	defer h.lock()()
	h.connCounter++
	h.conns[c] = h.connCounter
	log.Printf("Connection %d registered", h.conns[c])
}

// remove a connection from the hub and closes its channel.
func (h *Hub) remove(c *conn) {
	defer h.lock()()
	if id, exists := h.conns[c]; exists {
		delete(h.conns, c)
		close(c.output)
		log.Printf("Connection %d removed", id)
	} else {
		log.Println("Unknown connection")
	}
}

// lock locks the mutex and returns unlock as a function for use with defer.
func (h *Hub) lock() func() {
	h.mu.Lock()
	return func() { h.mu.Unlock() }
}
