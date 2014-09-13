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

	http.Handle("/", http.FileServer(http.Dir("./")))

	s := *self
	if s == "" {
		s = fmt.Sprintf("%s:%d", *addr, *port)
	}

	b := bridge.NewBridge(20, true) // Add flags.
	b.Start(s, *addr, *port)
}
