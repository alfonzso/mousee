package hid

import (
	// "fmt"
	"log"
	"math"
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

var _2On16 = int(math.Pow(float64(2), float64(16)))

func WheelMovement(mouseData uint32) uintptr {
	// sss = 4287102976
	sss := mouseData >> 16
	if int(sss)-_2On16 == -120 {
		return uintptr(common.WM_MOUSWHEELDOWN)
	} else if sss == 120 {
		return uintptr(common.WM_MOUSWHEELUP)
	} else {
		return 0
	}
	// fmt.Println(sss - p)
}

func KeyboardDefaultHookHandler(c chan<- types.KeyboardEvent) types.HOOKPROC {
	return func(code int32, wParam, lParam uintptr) uintptr {
		if lParam != 0 {
			keyBevt := types.KeyboardEvent{
				Message:         types.Message(wParam),
				KBDLLHOOKSTRUCT: *(*types.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam)),
			}
			// if keyBevt.VKCode.String() == "enter" {
			// log.Println(keyBevt)
			// log.Printf("%x %v",uint32(keyBevt.VKCode), keyBevt.VKCode)
			if keyBevt.VKCode == types.VK_ESCAPE {
				keyboardDebugMode += 1
			}
			c <- keyBevt
		}
		if keyboardDebugMode < 5 {
			return 1
		}
		if keyboardDebugMode == 6 {
			log.Println()
			log.Println("Keyboard enabled...")
			log.Println()
		}
		if keyboardDebugMode > 10 {
			keyboardDebugMode = 0
			return 1
		}
		return win32.CallNextHookEx(0, code, wParam, lParam)
	}
}

func MouseDefaultHookHandler(c chan<- types.MouseEvent) types.HOOKPROC {
	return func(code int32, wParam, lParam uintptr) uintptr {
		// mouseData := types.MouseEvent{}
		cont := true
		if lParam != 0 {
			mouseData := types.MouseEvent{
				Message:        types.Message(wParam),
				MSLLHOOKSTRUCT: *(*types.MSLLHOOKSTRUCT)(unsafe.Pointer(lParam)),
			}
			if mouseDebugMode < 5 {
				// select wParam {
				switch wParam {
				case uintptr(common.WM_MBUTTONDOWN), uintptr(common.WM_LBUTTONDOWN),
					uintptr(common.WM_RBUTTONDOWN), uintptr(common.WM_MOUSEWHEEL),
					uintptr(common.WM_MOUSEHWHEEL), uintptr(common.WM_XBUTTON4_5_DOWN),
					uintptr(common.WM_XBUTTON4_5_UP):
					// fmt.Println("Mouse click blocked!")
					cont = false
				}
			}

			switch wParam {
			case uintptr(common.WM_MOUSEWHEEL), uintptr(common.WM_MOUSEHWHEEL):
				mouseData.Message = types.Message(WheelMovement(mouseData.MouseData))
			}

			switch wParam {
			case uintptr(common.WM_XBUTTON4_5_DOWN):
				// log.Println(">>>>>>>>>>>", wParam, mouseData)
				// test := mouseData.MouseData >> 16
				// log.Println(test, mouseData.MouseData)
				// log.Printf(" %+v %d", mouseData, test)
				// mouseData.Message = types.Message(WheelMovement(mouseData.MouseData))
				if mouseData.MouseData >> 16 == 1{
					mouseData.Message = common.WM_MOUSE4BTN
				}
				if mouseData.MouseData >> 16 == 2{
					mouseData.Message = common.WM_MOUSE5BTN
				}
			}
			// log.Printf(" %+v ", mouseData)
			c <- mouseData
			// if mouseDebugMode < 5 && (wParam == uintptr(common.WM_LBUTTONDOWN) || wParam == uintptr(common.WM_RBUTTONDOWN)) {
			// 	log.Println("Mouse click blocked!")
			// 	cont = false
			// }
			if wParam == uintptr(common.WM_MBUTTONDOWN) {
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

		// log.Println(">>>>>>>>>>>", wParam, code, lParam)
		// switch wParam {
		// // case uintptr(common.WM_MBUTTONDOWN):
		// // case uintptr(common.WM_LBUTTONDOWN):
		// // case uintptr(common.WM_RBUTTONDOWN):
		// case uintptr(common.WM_MOUSEWHEEL), uintptr(common.WM_MOUSEHWHEEL):
		// 	// log.Println(">>>>>>>>>>>", wParam, mouseData)
		// 	log.Printf(" %+v ", mouseData)
		// default:
		// 	// freebsd, openbsd,
		// 	// plan9, windows...
		// 	// fmt.Printf("%s.\n", os)
		// 	// continue
		// }

		if !cont {
			return 1
		}

		// ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
		// return ret
		return win32.CallNextHookEx(0, code, wParam, lParam)
	}
}
