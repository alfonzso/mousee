package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	// "syscall"
	"time"
	"unsafe"

	"github.com/alfonzso/mousee/common"
	"github.com/alfonzso/mousee/server"

	// "github.com/moutend/go-hook/pkg/mouse"
	"github.com/moutend/go-hook/pkg/types"

	"github.com/moutend/go-hook/pkg/win32"

	"github.com/gorilla/websocket"
)

type HookHandler func(c chan<- types.MouseEvent) types.HOOKPROC

// var (
// 	user32                  = syscall.NewLazyDLL("user32.dll")
// 	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
// 	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
// 	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
// 	procGetMessage          = user32.NewProc("GetMessageW")

// 	hook syscall.Handle
// )

// const (
// 	WH_MOUSE_LL    = 14
// 	WM_LBUTTONDOWN = 0x0201
// 	WM_RBUTTONDOWN = 0x0204
// 	MK_MBUTTON     = 0x0207
// )

type POINT struct {
	X, Y int32
}

type MouseLLHookStruct struct {
	Point     POINT
	MouseData uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

var mouseDebugMode = 0

func MouseDefaultHookHandler(c chan<- types.MouseEvent) types.HOOKPROC {
	return func(code int32, wParam, lParam uintptr) uintptr {
		cont := true
		if lParam != 0 {
			// Intercept left and right mouse button down events
			c <- types.MouseEvent{
				Message:        types.Message(wParam),
				MSLLHOOKSTRUCT: *(*types.MSLLHOOKSTRUCT)(unsafe.Pointer(lParam)),
			}
			if mouseDebugMode < 5 && (wParam == uintptr(common.WM_LBUTTONDOWN) || wParam == uintptr(common.WM_RBUTTONDOWN)) {
				log.Println("Mouse click blocked!")
				// return 1 // Block the event
				cont = false
			}
			// if wParam == WM_LBUTTONDOWN || wParam == MK_MBUTTON {
			// 	log.Println("keeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
			// }
			// if wParam == MK_MBUTTON {
			if wParam == uintptr(common.WM_MBUTTON) {
				mouseDebugMode += 1
				if mouseDebugMode >= 5 {
					log.Println("Debug mode active for mouse", mouseDebugMode)
				}
				if mouseDebugMode > 10 {
					mouseDebugMode = 0
				}
				// log.Println("keeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
			}
			// log.Println(">>>>>>>>>>>", wParam)
		}

		// log.Println(">>>>>>>>>>>", wParam)

		if !cont {
			return 1
		}

		// ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
		// return ret
		return win32.CallNextHookEx(0, code, wParam, lParam)
	}
}

// func SendDataToClient(u *server.UdpConfig, signalChan chan os.Signal, mouseChan chan types.MouseEvent) error {
func SendDataToClient(ws *server.WSServer, signalChan chan os.Signal, mouseChan chan types.MouseEvent, keyboardChan chan types.KeyboardEvent) error {

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)

	// res := ws.IsClientConnected(signalChan)
	// if !res {
	// 	os.Exit(0)
	// }

	fmt.Println("start capturing mouse input")

	// if err := mouse.Install(nil, mouseChan); err != nil {
	// 	return err
	// }

	// defer mouse.Uninstall()
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
			fmt.Printf(">>k>> %+v \r", k)
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

// func Mouse() {
// func Mouse(fn HookHandler, c chan<- types.MouseEvent) error {
// 	// Set the low-level mouse hook

// 	if fn == nil {
// 		fn = DefaultHookHandler
// 	}

// 	hook, _, err := procSetWindowsHookEx.Call(
// 		uintptr(common.WH_MOUSE_LL),
// 		syscall.NewCallback(fn(c)),
// 		0,
// 		0,
// 	)
// 	if hook == 0 {
// 		log.Fatal("Failed to set mouse hook:", err)
// 	}
// 	defer procUnhookWindowsHookEx.Call(hook)

// 	// Wait for mouse events
// 	log.Println("Wait for mouse events")
// 	var msg struct {
// 		hwnd    uintptr
// 		message uint32
// 		wParam  uintptr
// 		lParam  uintptr
// 		time    uint32
// 		pt      POINT
// 	}
// 	for {
// 		procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
// 	}
// }
