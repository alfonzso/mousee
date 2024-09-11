package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"github.com/alfonzso/mousee/server"
	"github.com/moutend/go-hook/pkg/mouse"
	"github.com/moutend/go-hook/pkg/types"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage          = user32.NewProc("GetMessageW")

	hook syscall.Handle
)

const (
	WH_MOUSE_LL    = 14
	WM_LBUTTONDOWN = 0x0201
	WM_RBUTTONDOWN = 0x0204
)

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

func LowLevelMouseProc(nCode int32, wParam, lParam uintptr) uintptr {
	if nCode == 0 {
		// Intercept left and right mouse button down events
		if false && wParam == WM_LBUTTONDOWN || wParam == WM_RBUTTONDOWN {
			log.Println("Mouse click blocked!")
			return 1 // Block the event
		}
	}

	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}

func MousePosHook(u server.UdpConfig) error {

	mouseChan := make(chan types.MouseEvent, 100)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	res := u.IsClientConnected(signalChan)
	if !res {
		os.Exit(0)
	}

	fmt.Println("start capturing mouse input")

	if err := mouse.Install(nil, mouseChan); err != nil {
		return err
	}

	defer mouse.Uninstall()

	for {

		select {
		case <-time.After(5 * time.Minute):
			fmt.Println("Received timeout signal")
			return nil
		case <-signalChan:
			fmt.Println("Received shutdown signal")
			return nil
		case m := <-mouseChan:
			msg := fmt.Sprintf("Received %v {X:%v, Y:%v}\n", m.Message, m.X, m.Y)

			u.SendResponse(msg)
			continue
		}
	}
}

func Mouse() {
	// Set the low-level mouse hook
	hook, _, err := procSetWindowsHookEx.Call(
		WH_MOUSE_LL,
		syscall.NewCallback(LowLevelMouseProc),
		0,
		0,
	)
	if hook == 0 {
		log.Fatal("Failed to set mouse hook:", err)
	}
	defer procUnhookWindowsHookEx.Call(hook)

	// Wait for mouse events
	log.Println("Wait for mouse events")
	var msg struct {
		hwnd    uintptr
		message uint32
		wParam  uintptr
		lParam  uintptr
		time    uint32
		pt      POINT
	}
	for {
		procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
	}
}
