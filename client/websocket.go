package client

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/alfonzso/mousee/common"
	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"github.com/moutend/go-hook/pkg/types"
)

// var addr = flag.String("addr", "127.0.0.1:12345", "http service address")

var commonData common.CommonData

var remapAtoZKeys = map[types.VKCode]string{
	types.VK_A:      "a",
	types.VK_B:      "b",
	types.VK_C:      "c",
	types.VK_D:      "d",
	types.VK_E:      "e",
	types.VK_F:      "f",
	types.VK_G:      "g",
	types.VK_H:      "h",
	types.VK_I:      "i",
	types.VK_J:      "j",
	types.VK_K:      "k",
	types.VK_L:      "l",
	types.VK_M:      "m",
	types.VK_N:      "n",
	types.VK_O:      "o",
	types.VK_P:      "p",
	types.VK_Q:      "q",
	types.VK_R:      "r",
	types.VK_S:      "s",
	types.VK_T:      "t",
	types.VK_U:      "u",
	types.VK_V:      "v",
	types.VK_W:      "w",
	types.VK_X:      "x",
	types.VK_Y:      "y",
	types.VK_Z:      "z",
	types.VK_SHIFT:  "shift",
	types.VK_LSHIFT: "shift",
	types.VK_RSHIFT: "shift",
}

func decodeMouseData(message []byte) {
	if err := json.Unmarshal(message, &commonData); err != nil {
		panic(err)
	}

	// fmt.Printf("%+v   %d %d            \r", commonData, types.WM_KEYDOWN, commonData.Msg)

	if commonData.VKCode != 0 {

		fmt.Println(types.WM_KEYDOWN == types.Message(commonData.Msg))
		str, ok := remapAtoZKeys[commonData.VKCode]
		if !ok {
			str = commonData.VKCode.String()
		}
		fmt.Println(string(str))
		// robotgo.KeyDown("shift")
		if types.WM_KEYDOWN == types.Message(commonData.Msg) {
			robotgo.KeyToggle(str)
		} else {
			robotgo.KeyToggle(str, "up")
		}
		// robotgo.KeyUp("shift")
		return
	}

	if commonData.X != -1 && commonData.Y != -1 {
		robotgo.Move(int(commonData.X), int(commonData.Y))
	}

	switch commonData.Msg {
	case uintptr(common.WM_MBUTTONDOWN):
		robotgo.Toggle("center")
	case uintptr(common.WM_MBUTTONUP):
		robotgo.Toggle("center", "up")
	case uintptr(common.WM_LBUTTONUP):
		robotgo.Toggle("left", "up")
	case uintptr(common.WM_LBUTTONDOWN):
		robotgo.Toggle("left")
	case uintptr(common.WM_RBUTTONDOWN):
		robotgo.Toggle("right", "up")
	case uintptr(common.WM_RBUTTONUP):
		robotgo.Toggle("right")
	case uintptr(common.WM_MOUSEWHEEL):
		robotgo.ScrollDir(10, "down")
	case uintptr(common.WM_MOUSEHWHEEL):
		robotgo.ScrollDir(10, "up")
		log.Println("Mouse click blocked!")
		// cont = false
	}
	// if uintptr(common.WM_LBUTTONDOWN) == commonData.Msg {
	// 	robotgo.Toggle("left")
	// }

	// if uintptr(common.WM_LBUTTONUP) == commonData.Msg {
	// 	robotgo.Toggle("left", "up")
	// }

	// if uintptr(common.WM_RBUTTONDOWN) == commonData.Msg {
	// 	robotgo.Toggle("right")
	// }
	// if uintptr(common.WM_RBUTTONUP) == commonData.Msg {
	// 	robotgo.Toggle("right", "up")
	// }
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
