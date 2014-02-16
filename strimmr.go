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
var index = template.Must(template.ParseFiles("index.html"))

func main() {
	hub := NewHub()

	flag.Parse()
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/socket", hub.HandleConnection)

	go func() {
		for {
			hub.Broadcast([]byte("One Two Three Four"))
			time.Sleep(time.Second * 2)
		}
	}()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
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
