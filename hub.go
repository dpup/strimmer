package main

import (
	"log"
	"sync"
	"time"
)

// Hub keeps a list of connections and broadcasts data to them.
type Hub struct {
	conns map[*Conn]time.Time
	mu    sync.Mutex // protects conns
}

func NewHub() *Hub {
	return &Hub{conns: make(map[*Conn]time.Time)}
}

// Register a new  connection.
func (h *Hub) Register(c *Conn) {
	defer h.lock()()
	h.conns[c] = time.Now()
	log.Println("New connection registered")
}

// Unregister a connection from the hub and closes its channel.
func (h *Hub) Unregister(c *Conn) {
	defer h.lock()()
	delete(h.conns, c)
	close(c.output)
	log.Println("Connection removed")
}

func (h *Hub) Broadcast(data []byte) {
	defer h.lock()()
	for c, started := range h.conns {
		select {
		case c.output <- data:
			log.Println("Sending message :", started)
		default:
			log.Println("Connection closed :", started)
			delete(h.conns, c)
			close(c.output)
		}
	}
}

// lock locks the mutex and returns unlock as a function.
func (h *Hub) lock() func() {
	h.mu.Lock()
	return func() { h.mu.Unlock() }
}
