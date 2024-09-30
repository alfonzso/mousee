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
	types.VK_A:        "a",
	types.VK_B:        "b",
	types.VK_C:        "c",
	types.VK_D:        "d",
	types.VK_E:        "e",
	types.VK_F:        "f",
	types.VK_G:        "g",
	types.VK_H:        "h",
	types.VK_I:        "i",
	types.VK_J:        "j",
	types.VK_K:        "k",
	types.VK_L:        "l",
	types.VK_M:        "m",
	types.VK_N:        "n",
	types.VK_O:        "o",
	types.VK_P:        "p",
	types.VK_Q:        "q",
	types.VK_R:        "r",
	types.VK_S:        "s",
	types.VK_T:        "t",
	types.VK_U:        "u",
	types.VK_V:        "v",
	types.VK_W:        "w",
	types.VK_X:        "x",
	types.VK_Y:        "y",
	types.VK_Z:        "z",
	types.VK_SHIFT:    "shift",
	types.VK_LSHIFT:   "lshift",
	types.VK_RSHIFT:   "rshift",
	types.VK_RETURN:   "enter",
	types.VK_CONTROL:  "ctrl",
	types.VK_LCONTROL: "lctrl",
	types.VK_RCONTROL: "rctrl",
	types.VK_MENU:     "alt",
	types.VK_LMENU:    "lalt",
	types.VK_RMENU:    "ralt",
	types.VK_TAB:      "tab",
	types.VK_CAPITAL:  "capslock",
	types.VK_SPACE:    "space",
	types.VK_INSERT:   "inert",
	types.VK_ESCAPE:   "esc",
	types.VK_UP:       "up",
	types.VK_DOWN:     "down",
	types.VK_LEFT:     "left",
	types.VK_RIGHT:    "right",
	types.VK_HOME:     "home",
	types.VK_DELETE:   "del",
	types.VK_END:      "end",

	types.VK_OEM_COMMA:  ",",
	types.VK_OEM_MINUS:  "-",
	types.VK_OEM_PLUS:   "+",
	types.VK_OEM_PERIOD: ".",

	types.VK_OEM_1: ";",  //     For the US standard keyboard, the ';:' key
	types.VK_OEM_2: "/",  //     For the US standard keyboard, the '/?' key
	types.VK_OEM_3: "`",  //      For the US standard keyboard, the '`~' key
	types.VK_OEM_4: "[",  //      For the US standard keyboard, the '[{' key
	types.VK_OEM_5: "\\", //      For the US standard keyboard, the '\|' key
	types.VK_OEM_6: "]",  //     For the US standard keyboard, the ']}' key
	types.VK_OEM_7: "'",  //      For the US standard keyboard, the 'single-quote/double-quote' key
	// VK_OEM_8
}

// var keyAsUint32 string

func decodeMouseData(message []byte) {
	if err := json.Unmarshal(message, &commonData); err != nil {
		panic(err)
	}

	// fmt.Printf("%+v   %d %d            \r", commonData, types.WM_KEYDOWN, commonData.Msg)

	if commonData.VKCode != 0 {
		key, ok := remapAtoZKeys[commonData.VKCode]
		if !ok {
			key = string(rune(uint32(commonData.VKCode)))
		}
		if types.WM_KEYDOWN == types.Message(commonData.Msg) || types.WM_SYSKEYDOWN == types.Message(commonData.Msg) {
			robotgo.KeyToggle(key)
			// fmt.Printf(">d> %s\n", key)
		} else {
			robotgo.KeyToggle(key, "up")
			// fmt.Printf(">u> %s\n", key)
		}
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
	case uintptr(common.WM_MOUSWHEELDOWN):
		robotgo.ScrollDir(5, "down")
	case uintptr(common.WM_MOUSWHEELUP):
		robotgo.ScrollDir(5, "up")
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
