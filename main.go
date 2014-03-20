// ipecho project main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort = "8080"
	defaultHost = "0.0.0.0"
)

func render(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/plain")
		address, _, _ := net.SplitHostPort(r.RemoteAddr)
		io.WriteString(w, address+"\n")
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	var listenPort string
	var listenHost string
	var help bool

	flag.StringVar(&listenPort, "port", defaultPort, "TCP port to listen on")
	flag.StringVar(&listenHost, "host", defaultHost, "TCP address to listen on")
	flag.BoolVar(&help, "help", false, "show usage message")
	flag.Parse()

	listenAddress := net.JoinHostPort(listenHost, listenPort)
	server := &http.Server{
		Addr:        listenAddress,
		ReadTimeout: 1 * time.Second,
	}
	http.HandleFunc("/", render)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem listening: %s\n", err.Error())
		os.Exit(1)
	}

}
