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

	"github.com/ISBachvarov21/Rune-GO-test/routes"
	"github.com/ISBachvarov21/Rune-GO/runehttp"
)

var portChan chan string
var port string

type Serializable struct {
	S string `json:"s"`
}

func Root(c *runehttp.Context) (runehttp.Response, error) {
	return runehttp.Text("ok "), nil
}

func SerializableTest(s Serializable) (runehttp.Response, error) {
	return runehttp.Text(s.S), nil
}

func server(srv *http.Server) {
	var handler = runehttp.NewHTTPHandler()

	err := routes.RegisterRoutes(handler)
	if err != nil {
		panic(err)
	}

	srv.Handler = handler
	srv.Addr = ":" + port
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe failed: %v", err)
	}
}

func main() {
	// todo: wait for port and report start only if ran with specific flag (-dev)
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
