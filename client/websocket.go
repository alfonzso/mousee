package client

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// var addr = flag.String("addr", "127.0.0.1:12345", "http service address")

func WsClientMode() {
	infoLogger := log.New(os.Stdout, "INFO: ", 0)

	infoLogger.Println("Client mode active ...")

	addr := flag.String("cliAddr", "localhost:5555", "http service address")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/client"}

	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	log.Println("connected")

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	log.Println("gooooooooo")
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	log.Printf("connecting to %s", "writing....")
	c.WriteMessage(websocket.TextMessage, []byte("Hi UDP Server, How are you doing?"))

	for {
		select {
		case <-done:
			return

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
