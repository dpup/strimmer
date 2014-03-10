// Copyright 2014 Daniel Pupius
// Based on tutorial at http://gary.burd.info/go-websocket-chat

package main

import (
	"container/ring"
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
	conns       map[*Conn]int
	mu          sync.Mutex // protects conns
	connCounter int
	history     *ring.Ring
}

// Creates a new hub for managing websocket connections.
func NewHub(historySize int) *Hub {
	return &Hub{conns: make(map[*Conn]int), history: ring.New(historySize)}
}

// Handle connection responds to a HTTP request and attempts to upgrade the
// request into a websocket connection.
func (h *Hub) HandleConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method nod allowed", http.StatusMethodNotAllowed)
		return
	}
	// NOTE(dan): Allow cross-origin web sockets.
	// if r.Header.Get("Origin") != "http://"+r.Host {
	// 	http.Error(w, "Origin not allowed", http.StatusForbidden)
	// 	return
	// }

	socket, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		log.Println("Not a websocket handshake")
		http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		return
	} else if err != nil {
		log.Println("Error upgrading socket", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn := &Conn{output: make(chan []byte, 256), socket: socket, created: time.Now()}

	// Send the history to clients when they connect.
	h.history.Do(func(v interface{}) {
		if v != nil {
			if data, ok := v.([]byte); ok {
				conn.output <- data
			}
		}
	})

	h.add(conn)
	defer func() { h.remove(conn) }()
	conn.WritePump()
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
	if h.history.Len() > 0 {
		h.history.Value = data
		h.history = h.history.Next()
	}
}

// add a new  connection.
func (h *Hub) add(c *Conn) {
	defer h.lock()()
	h.connCounter++
	h.conns[c] = h.connCounter
	log.Printf("Connection %d registered", h.conns[c])
}

// remove a connection from the hub and closes its channel.
func (h *Hub) remove(c *Conn) {
	defer h.lock()()
	if id, exists := h.conns[c]; exists {
		delete(h.conns, c)
		close(c.output)
		log.Printf("Connection %d removed", id)
	} else {
		log.Println("Unknown collection")
	}
}

// lock locks the mutex and returns unlock as a function for use with defer.
func (h *Hub) lock() func() {
	h.mu.Lock()
	return func() { h.mu.Unlock() }
}
