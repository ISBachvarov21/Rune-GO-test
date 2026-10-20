package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var portChan chan string
var port string

func server(srv *http.Server) {
	var handler = &httphandler.HTTPHandler{}

	handler.AddRoute("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("ivan e gotin!!")
	})

	//handler.AddRoute("POST", "/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("POST")
	//})

	srv.Handler = handler
	srv.Addr = ":" + port
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe failed: %v", err)
	}
}

func main() {
	shutdown := make(chan struct{})
	portChan = make(chan string)
	var srv http.Server

	go func() {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			text := scanner.Text()
			if text == "shutdown\n" || text == "shutdown" {
				close(shutdown)
			}
			if strings.Contains(text, "port") {
				portChan <- strings.Split(text, ":")[1]
			}
		}
	}()

	port = <-portChan

	go server(&srv)
	go func() {
		client := &http.Client{}
		defer client.CloseIdleConnections()
		for {
			_, err := client.Head("http://localhost:" + port + "/readinesscheck")
			if err == nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println("running")
	}()

	for {
		select {
		case <-shutdown:
			srv.Shutdown(context.Background())
			return
		}
	}
}
