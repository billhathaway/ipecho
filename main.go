// ipecho project main.go
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	defaultPort          = "8080"
	defaultHost          = "0.0.0.0"
	mainPageTemplateFile = "templates/mainPage.tmpl"
)

type (
	clientInfo struct {
		IP       string
		Hostname string
	}
)

var (
	mainPageTemplate = template.Must(template.ParseFiles(mainPageTemplateFile))
)

func log(message string) {
	fmt.Println(time.Now().Format(time.RFC3339) + " " + message)
}

func resolve(remoteAddress string) (*clientInfo, error) {
	client := new(clientInfo)
	client.Hostname = "UNRESOLVED"
	var err error
	client.IP, _, err = net.SplitHostPort(remoteAddress)
	if err != nil {
		return client, err
	}
	names, err := net.LookupAddr(client.IP)
	if err != nil {
		return client, err
	}
	if len(names) > 0 {
		client.Hostname = names[0]
	}
	return client, nil
}

func logRequest(path string, client *clientInfo, responseTime time.Duration, resolveErr error, writeErr error) {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("path=%s ip=%s name=%s time=%f", path, client.IP, client.Hostname, responseTime.Seconds()))
	if resolveErr != nil {
		buf.WriteString(" resolveErr=" + resolveErr.Error())
	}
	if writeErr != nil {
		buf.WriteString(" writeErr=" + writeErr.Error())
	}
	log(buf.String())
}

func render(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	client, resolveErr := resolve(r.RemoteAddr)
	var writeErr error
	switch r.URL.Path {
	default: // for any page that we don't handle explicitly, render the main page
		fallthrough
	case "/":
		writeErr = mainPageTemplate.Execute(w, client)
	case "/text":
		_, writeErr = io.WriteString(w, client.IP)
	case "/json":
		output, _ := json.Marshal(client)
		_, writeErr = w.Write(output)
	case "/favicon.ico":
		http.ServeFile(w, r, "public/icons/favicon.ico")
		return // don't bother logging these
	}

	logRequest(r.URL.Path, client, time.Since(startTime), resolveErr, writeErr)
}

func main() {
	var listenPort string
	var listenHost string

	flag.StringVar(&listenPort, "port", defaultPort, "TCP port to listen on")
	flag.StringVar(&listenHost, "host", defaultHost, "TCP address to listen on")
	flag.Parse()

	listenAddress := net.JoinHostPort(listenHost, listenPort)
	server := &http.Server{
		Addr:        listenAddress,
		ReadTimeout: 10 * time.Second,
	}
	http.HandleFunc("/", render)
	log(fmt.Sprintf("event=start address=%s", listenAddress))
	server.ListenAndServe()

}
