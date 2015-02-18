// Copyright 2014 Daniel Pupius

package main

import (
	"container/ring"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/dpup/strimmer/bridge"
)

var self = flag.String("self", "", "hostname:port that remote clients should connect to, even if behind a proxy")
var port = flag.Int("port", 3100, "Port to listen on")
var addr = flag.String("addr", "", "IP address the server should listen on")

func main() {
	flag.Parse()

	lr := &logRecorder{logs: ring.New(1000)}
	log.SetOutput(lr)

	http.HandleFunc("/debug/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		lr.Dump(w)
	})
	http.Handle("/", http.FileServer(http.Dir("./")))

	s := *self
	if s == "" {
		s = fmt.Sprintf("%s:%d", *addr, *port)
	}

	b := bridge.NewBridge(20, true) // Add flags.
	b.Start(s, *addr, *port)
}

type logRecorder struct {
	logs *ring.Ring
	mu   sync.RWMutex
}

func (r *logRecorder) Dump(w io.Writer) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.logs.Do(func(line interface{}) {
		if line != nil {
			w.Write(line.([]byte))
		}
	})
}

func (r *logRecorder) Write(p []byte) (n int, err error) {
	cp := make([]byte, len(p), len(p))
	copy(cp, p)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs.Value = cp
	r.logs = r.logs.Next()
	return os.Stderr.Write(p)
}
