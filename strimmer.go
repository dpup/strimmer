// Copyright 2014 Daniel Pupius

package main

import (
	"flag"
	"html/template"
	"net/http"

	"github.com/dpup/strimmer/bridge"
)

var port = flag.Int("port", 3100, "Port to listen on")
var host = flag.String("host", "", "Host or IP of the server")

func main() {
	flag.Parse()

	http.HandleFunc("/", handleHome)

	b := bridge.NewBridge(20, true) // Add flags.
	b.Start(*host, *port)
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
