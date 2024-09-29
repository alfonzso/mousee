package common

import (
	"os"
)

var AppVersion string
var AppName string

type MouseData struct {
	Msg uintptr
	X   int32 `json:"X"`
	Y   int32 `json:"Y"`
}
type KeyBoardData struct {
	Msg uintptr
}

type UpdateData struct {
	AppName    string
	AppVersion string
	AppCrc32   uint32
}

// HOOKPROC represents HOOKPROC callback function type.
//
// For more details, see the MSDN.
//
// https://docs.microsoft.com/en-us/windows/win32/winmsg/using-hooks
type HOOKPROC func(code int32, wParam, lParam uintptr) uintptr

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func BeginUpdate() string {
	begin := []rune{'B', 'E', 'G', 'I', 'N', '_', 'U', 'P', 'D', 'A', 'T', 'E'}
	result := ""
	for _, v := range begin {
		result += string(v)
	}
	return result
}

func EndUpdate() string {
	begin := []rune{'E', 'N', 'D', '_', 'U', 'P', 'D', 'A', 'T', 'E'}
	result := ""
	for _, v := range begin {
		result += string(v)
	}
	return result
}

func UpdateFile(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	return f
}
