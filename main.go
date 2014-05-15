// ipecho project main.go
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort = "8080"
	defaultHost = "0.0.0.0"
)

var (
	quiet = false
)

func render(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")
	address, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Fprintln(w, address)
	if !quiet {
		fmt.Printf("%s ip=%s request=%s\n", time.Now().Format(time.RFC3339), address, r.URL.Path)
	}

}

func main() {
	var listenPort string
	var listenHost string
	var help bool

	flag.StringVar(&listenPort, "port", defaultPort, "TCP port to listen on")
	flag.StringVar(&listenHost, "host", defaultHost, "TCP address to listen on")
	flag.BoolVar(&quiet, "quiet", quiet, "disable request logging")
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
