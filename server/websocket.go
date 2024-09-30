package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alfonzso/mousee/common"
	"github.com/gorilla/websocket"
	"github.com/moutend/go-hook/pkg/types"
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

func (ws *WSServer) SendDataToClient(signalChan chan os.Signal, mouseChan chan types.MouseEvent, keyboardChan chan types.KeyboardEvent) error {

	fmt.Println("start capturing mouse input")

	var isClientCOnnected bool
	for {
		// time.Sleep(100 * time.Millisecond)
		select {
		case <-time.After(5 * time.Minute):
			fmt.Println("Received timeout signal")
			return nil
		case <-signalChan:
			fmt.Println("Received shutdown signal")
			return nil
		case isClientCOnnected = <-ws.ClientConnected:
			continue
		case k := <-keyboardChan:
			// fmt.Printf(">>k>> %+v \r", k)
			if isClientCOnnected {
				// b, err := json.Marshal(common.MouseData{X: k.X, Y: m.Y, Msg: uintptr(m.Message)})
				f := common.KeyBoardData{VKCode: k.VKCode, X: -1, Y: -1, Msg: uintptr(k.Message)}
				// f.X
				// b, err := json.Marshal(common.KeyBoardData{X: 0, Y: 0, Msg: uintptr(k.Message), VKCode: k.VKCode})
				b, err := json.Marshal(f)
				if err == nil {
					// ws.SendResponse(string(b) + "\n")
					for conn := range ws.Clients {
						conn.WriteMessage(websocket.TextMessage, b)
					}
				}
			}
			continue
		case m := <-mouseChan:
			// msg := fmt.Sprintf("Received %v {X:%v, Y:%v}\n", m.Message, m.X, m.Y)
			// msg := fmt.Sprintf("%v %v", m.X, m.Y)
			// md := common.MouseData{m.X, m.Y}
			// WM_LBUTTONDOWN
			// types.Message(wParam)
			// fmt.Printf("%v \r", string(b))
			// if false {
			if isClientCOnnected {
				b, err := json.Marshal(common.MouseData{VKCode: 0, X: m.X, Y: m.Y, Msg: uintptr(m.Message)})
				if err == nil {
					// ws.SendResponse(string(b) + "\n")
					for conn := range ws.Clients {
						conn.WriteMessage(websocket.TextMessage, b)
					}
				}
			}
			continue
		}
	}
}
