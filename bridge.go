// Copyright 2014 Daniel Pupius
// Based on tutorial at http://gary.burd.info/go-websocket-chat

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dpup/gohubbub"
	"github.com/gorilla/websocket"
)

// Bridge keeps broadcasts messages recieved from PuSH to websocket clients.
type Bridge struct {
	pushClient  *gohubbub.Client
	conns       map[*conn]int
	mu          sync.Mutex // protects conns
	connCounter int
	// history          *ring.Ring
	allowCrossOrigin bool
}

// Creates a new hub for managing websocket connections.
func NewBridge(historySize int, allowCrossOrigin bool) *Bridge {
	return &Bridge{
		conns: make(map[*conn]int),
		// history:          ring.New(historySize),
		allowCrossOrigin: allowCrossOrigin,
	}
}

// Start registers handles on DefaultServeMux and starts up a HTTP server using
// ListenAndServe.
func (b *Bridge) Start(host string, port int) {
	log.Printf("Setting up server on %s:%d", host, port)
	b.pushClient = gohubbub.NewClient(host, port, "strimmr")
	b.pushClient.RegisterHandler(http.DefaultServeMux)
	go b.pushClient.Start()
	http.HandleFunc("/bridge", b.HandleConnection)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func (b *Bridge) Shutdown() {
	log.Println("Shutting down...")
	for c, _ := range b.conns {
		b.remove(c) // Triggers remove.
	}
}

// HandleConnection responds to a HTTP request and attempts to upgrade the
// request into a websocket connection.
func (b *Bridge) HandleConnection(w http.ResponseWriter, r *http.Request) {
	feedUrl := r.URL.Query().Get("feed")

	log.Println(r.URL, feedUrl)

	if !b.pushClient.HasSubscription(feedUrl) {
		err := b.pushClient.DiscoverAndSubscribe(feedUrl, func(contentType string, body []byte) {
			b.broadcast(feedUrl, body)
		})
		if err != nil {
			log.Println("Error subscribing to hub,", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("Subscription already exists for", feedUrl)
	}

	socket, err, statusCode := b.upgrade(w, r)
	if err != nil {
		log.Println("Error handling conncetion,", err)
		http.Error(w, err.Error(), statusCode)
		return
	}

	conn := &conn{
		feedUrl: feedUrl,
		output:  make(chan []byte, 256),
		socket:  socket,
		created: time.Now(),
	}

	// Send recent history to clients when they connect.
	// b.history.Do(func(v interface{}) {
	// 	if v != nil {
	// 		if data, ok := v.([]byte); ok {
	// 			conn.output <- data
	// 		}
	// 	}
	// })

	b.add(conn)

	defer func() { b.remove(conn) }()
	conn.wait()
}

func (b *Bridge) upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error, int) {
	if r.Method != "GET" {
		return nil, errors.New("method not allowed"), http.StatusMethodNotAllowed
	}

	if !b.allowCrossOrigin && r.Header.Get("Origin") != "http://"+r.Host {
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

// broadcast sends data to each matching connection.
func (b *Bridge) broadcast(feedUrl string, data []byte) {
	// TODO: Remove this lock, need to copy conns and then have
	// a separate cleanup pass for dead subscriptions.
	defer b.lock()()

	feed, err := XMLToFeed(data)
	if err != nil {
		log.Println("Error processing hub response", err)
	} else {

		jsonFeed, _ := json.Marshal(feed)

		count := 0
		for c, id := range b.conns {
			// TODO(dan): Use a map to look up connections that match a feed.
			if c.feedUrl == feedUrl {
				select {
				case c.output <- jsonFeed:
					count++
				default:
					log.Println("Connection closed :", id)
					delete(b.conns, c)
					close(c.output)
				}
			}
		}

		for _, entry := range feed.Entries {
			log.Printf("Broadcasted to %d clients : %s", count, entry.Title)
		}
	}

	b.checkFeedHasClients(feedUrl)

	// Store the history.
	// b.history.Value = data
	// b.history = b.history.Next()
}

// add a new connection.
func (b *Bridge) add(c *conn) {
	defer b.lock()()
	b.connCounter++
	b.conns[c] = b.connCounter
	log.Printf("Connection %d registered", b.conns[c])
}

// remove a connection from the hub and closes its channel.
func (b *Bridge) remove(c *conn) {
	defer b.lock()()
	if id, exists := b.conns[c]; exists {
		delete(b.conns, c)
		close(c.output)
		log.Printf("Connection %d (%s) removed", id, c.feedUrl)
	}
	b.checkFeedHasClients(c.feedUrl)
}

// checkFeedHasClients will unsubscribe from a feed if all the connections have
// been removed. Should be called within a lock.
func (b *Bridge) checkFeedHasClients(feedUrl string) {
	for c, _ := range b.conns {
		if c.feedUrl == feedUrl {
			return
		}
	}
	if b.pushClient.HasSubscription(feedUrl) {
		b.pushClient.Unsubscribe(feedUrl)
	}
}

// lock locks the mutex and returns unlock as a function for use with defer.
func (b *Bridge) lock() func() {
	b.mu.Lock()
	return func() { b.mu.Unlock() }
}
