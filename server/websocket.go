package server

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type WsConfig struct {
	Addr            *net.UDPAddr
	Conn            *net.UDPConn
	Remoteaddr      *net.UDPAddr
	ClientConnected chan bool
}

//	var Upgrader = websocket.Upgrader{
//		CheckOrigin: func(r *http.Request) bool {
//			return true // Accepting all requests
//		},
//	}
var Upgrader = websocket.Upgrader{} // use default options

type WSServer struct {
	Clients map[*websocket.Conn]bool
	// handleMessage   func(message []byte) // New message handler
	ClientConnected chan bool
}

// func (ws *WsConfig) ServeWS(handleMessage func(message []byte)) *Server {
var flagAddr = flag.String("wsAddr", "192.168.1.100:5555", "http service address")

// func home(w http.ResponseWriter, r *http.Request) {
// 	// homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
// 	w.Write([]byte("Here is a string...."))
// }

type REQ func(http.ResponseWriter, *http.Request)

func ServeWS() *WSServer {
	flag.Parse()
	server := WSServer{
		make(map[*websocket.Conn]bool),
		make(chan bool),
	}

	http.HandleFunc("/client", server.Client())
	go func() {
		fmt.Println(">> WSServer listening at", *flagAddr)
		errr := http.ListenAndServe(*flagAddr, nil)
		fmt.Println(">> ", errr)
	}()

	// server.IsClientConnected(nil)

	return &server
}

// The function will block the runtime until someone connects to us.
// func (u WsConfig) IsClientConnected(signalChan chan os.Signal) bool {
func (serv *WSServer) IsClientConnected(signalChan chan os.Signal) bool {
	fmt.Println(">> IsClientConnected")
	for {
		// time.Sleep(100 * time.Millisecond)
		select {
		case connected := <-serv.ClientConnected:
			if connected {
				fmt.Println(">> ClientConnected")
				return true
			}

		case <-signalChan:
			fmt.Println(">> App shutting down")
			return false
		}
	}
}

// func (serv *WSServer) Client(w http.ResponseWriter, r *http.Request) {
func (serv *WSServer) Client() REQ {
	return func(w http.ResponseWriter, r *http.Request) {
		serv.ClientConnected <- false
		connection, _ := Upgrader.Upgrade(w, r, nil)

		serv.Clients[connection] = true

		mt, message, err := connection.ReadMessage()

		if err != nil {
			fmt.Printf("Some error  %v", err)
			panic(err)
		}

		// fmt.Printf("Read a message from %s \n", message[:mt])
		fmt.Printf("Read a message from %s %d\n", message, mt)

		// serv.Remoteaddr = remoteaddr

		fmt.Println("Starting ...", connection.RemoteAddr())
		serv.ClientConnected <- true

	}
}
