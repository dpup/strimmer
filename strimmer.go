// Copyright 2014 Daniel Pupius

package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var port = flag.Int("port", 3100, "Port to serve on")
var host = flag.String("host", "", "Host or IP to serve from")
var index = template.Must(template.ParseFiles("index.html"))

func main() {
	flag.Parse()

	http.HandleFunc("/", handleHome)

	bridge := NewBridge(20, true) // Add flags.
	go bridge.Start(*host, *port)

	// Wait for user input before shutting down.
	log.Println("Press Enter for graceful shutdown...")

	var input string
	fmt.Scanln(&input)

	bridge.Shutdown()
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
