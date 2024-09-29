package main

import (
	"log"
	"unsafe"

	"github.com/alfonzso/mousee/common"
	"github.com/moutend/go-hook/pkg/types"
	"github.com/moutend/go-hook/pkg/win32"
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

// type POINT struct {
// 	X, Y int32
// }
//
// type MouseLLHookStruct struct {
// 	Point     POINT
// 	MouseData uint32
// 	Flags     uint32
// 	Time      uint32
// 	ExtraInfo uintptr
// }

var mouseDebugMode = 0
var keyboardDebugMode = 0

func KeyboardDefaultHookHandler(c chan<- types.KeyboardEvent) types.HOOKPROC {
	return func(code int32, wParam, lParam uintptr) uintptr {
		if lParam != 0 {
			keyBevt := types.KeyboardEvent{
				Message:         types.Message(wParam),
				KBDLLHOOKSTRUCT: *(*types.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam)),
			}
			// if keyBevt.VKCode.String() == "enter" {
			log.Println(keyBevt)
			if keyBevt.VKCode == types.VK_ESCAPE {
				keyboardDebugMode += 1
			}
			c <- keyBevt
		}
		if keyboardDebugMode > 10 {
			keyboardDebugMode = 0
		}
		if keyboardDebugMode > 5 {
			return 1
		}
		return win32.CallNextHookEx(0, code, wParam, lParam)
	}
}

func MouseDefaultHookHandler(c chan<- types.MouseEvent) types.HOOKPROC {
	return func(code int32, wParam, lParam uintptr) uintptr {
		cont := true
		if lParam != 0 {
			c <- types.MouseEvent{
				Message:        types.Message(wParam),
				MSLLHOOKSTRUCT: *(*types.MSLLHOOKSTRUCT)(unsafe.Pointer(lParam)),
			}
			if mouseDebugMode < 5 && (wParam == uintptr(common.WM_LBUTTONDOWN) || wParam == uintptr(common.WM_RBUTTONDOWN)) {
				log.Println("Mouse click blocked!")
				cont = false
			}
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
