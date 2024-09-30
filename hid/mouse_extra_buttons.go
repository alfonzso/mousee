package hid

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32        = windows.NewLazySystemDLL("user32.dll")
	procSendInput = user32.NewProc("SendInput")
)

const (
	MOUSEEVENTF_XDOWN = 0x0081
	MOUSEEVENTF_XUP   = 0x0100
	XBUTTON1          = 0x0001
	XBUTTON2          = 0x0002
)

type MOUSEINPUT struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type INPUT struct {
	Type uint32
	MI   MOUSEINPUT
}

func SimulateXButtonPress(xButton uint32) uintptr {
	var inputs [2]INPUT

	// Press XButton
	inputs[0].Type = 0 // INPUT_MOUSE
	inputs[0].MI.DwFlags = MOUSEEVENTF_XDOWN
	inputs[0].MI.MouseData = xButton

	// Release XButton
	inputs[1].Type = 0 // INPUT_MOUSE
	inputs[1].MI.DwFlags = MOUSEEVENTF_XUP
	inputs[1].MI.MouseData = xButton

	nInputs := uint32(len(inputs))
	ret, _, _ := procSendInput.Call(
		uintptr(nInputs),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)
	return ret
}
