// Copyright 2014 Daniel Pupius

package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/dpup/strimmer/bridge"
)

var port = flag.Int("port", 3100, "Port to listen on")
var addr = flag.String("addr", "", "Hostname or IP of the server")
var self = flag.String("self", "", "Address that remote clients should connect to")

func main() {
	flag.Parse()

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/styles.css", handleStyles)

	s := *self
	if s == "" {
		s = fmt.Sprintf("%s:%d", *addr, *port)
	}

	b := bridge.NewBridge(20, true) // Add flags.
	b.Start(s, *addr, *port)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleStyles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "styles.css")
}
