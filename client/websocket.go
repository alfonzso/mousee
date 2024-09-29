package client

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/alfonzso/mousee/common"
	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
)

// var addr = flag.String("addr", "127.0.0.1:12345", "http service address")

var mouseData common.MouseData

func decodeMouseData(message []byte) {
	if err := json.Unmarshal(message, &mouseData); err != nil {
		panic(err)
	}

	// fmt.Println("mouseData", string(message))
	robotgo.Move(int(mouseData.X), int(mouseData.Y))
	
	if uintptr(common.WM_LBUTTONDOWN) == mouseData.Msg {
		// fmt.Println(">>> left")
		robotgo.Toggle("left")
		// robotgo.Click("left")
	}

	if uintptr(common.WM_LBUTTONUP) == mouseData.Msg {
		// fmt.Println(">>> left")
		robotgo.Toggle("left", "up")
		// robotgo.Click("left")
	}

	if uintptr(common.WM_RBUTTONDOWN) == mouseData.Msg {
		// fmt.Println(">>> right")
		// robotgo.Click("right")
		robotgo.Toggle("right")
	}
	if uintptr(common.WM_RBUTTONUP) == mouseData.Msg {
		// fmt.Println(">>> right")
		// robotgo.Click("right")
		robotgo.Toggle("right", "up")
	}
}

func WsClientMode() {
	infoLogger := log.New(os.Stdout, "INFO: ", 0)

	infoLogger.Println("Client mode active ...")

	addr := flag.String("cliAddr", "192.168.1.100:5555", "http service address")

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
			decodeMouseData(message)
			// log.Printf("recv: %s", message)
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
