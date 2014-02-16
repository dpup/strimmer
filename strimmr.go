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

var hub = NewHub()

func main() {
	flag.Parse()
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/socket", handleSocket)
	go sendData()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func sendData() {
	for {
		hub.Broadcast([]byte("One Two Three Four"))
		time.Sleep(time.Second * 2)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != "GET" {
		methodNotAllowed(w)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	index.Execute(w, nil)
}

func handleSocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		methodNotAllowed(w)
		return
	}
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", http.StatusForbidden)
		return
	}
	conn := CreateConn(w, r)
	if conn != nil {
		hub.Register(conn)
		defer func() { hub.Unregister(conn) }()
		conn.WritePump()
	}
}

func methodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "Method nod allowed", http.StatusMethodNotAllowed)
}
