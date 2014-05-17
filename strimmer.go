// Copyright 2014 Daniel Pupius

package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"

	"github.com/dpup/strimmer/bridge"
)

var port = flag.Int("port", 3100, "Port to listen on")
var addr = flag.String("addr", "", "Hostname or IP of the server")
var self = flag.String("self", "", "Address that remote clients should connect to")

func main() {
	flag.Parse()

	http.HandleFunc("/", handleHome)

	s := *self
	if s == "" {
		s = fmt.Sprintf("%s:%d", *addr, *port)
	}

	b := bridge.NewBridge(20, true) // Add flags.
	b.Start(s, *addr, *port)
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
	index := template.Must(template.ParseFiles("index.html"))
	index.Execute(w, nil)
}
