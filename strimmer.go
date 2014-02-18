// Copyright 2014 Daniel Pupius

package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/dpup/gohubbub"
)

var port = flag.Int("port", 3100, "Port to serve on")
var host = flag.String("host", "", "Host or IP to serve from")
var index = template.Must(template.ParseFiles("index.html"))

func main() {
	flag.Parse()

	http.HandleFunc("/", handleHome)

	hub := NewHub(20)
	http.HandleFunc("/socket", hub.HandleConnection)

	push := gohubbub.NewClient("http://medium.superfeedr.com", *host, *port, "Strimmr")
	push.RegisterHandler(http.DefaultServeMux)
	push.Subscribe("https://medium.com/feed/latest", func(contentType string, body []byte) {
		json, err := XMLFeedToJSON(body)
		if err != nil {
			log.Println("Error processing hub response", err)
		} else {
			hub.Broadcast(json)
		}
	})

	// Start the default server.
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	}()

	// Start the PuSH client.
	go push.Run()

	// Wait for user input before shutting down.

	log.Println("Press Enter for graceful shutdown...")

	var input string
	fmt.Scanln(&input)

	push.Unsubscribe("https://medium.com/feed/latest")

	time.Sleep(time.Second * 5)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method nod allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	index.Execute(w, nil)
}
