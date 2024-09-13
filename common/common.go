package common

type MouseData struct {
	Msg uintptr
	X   int32 `json:"X"`
	Y   int32 `json:"Y"`
}

// HOOKPROC represents HOOKPROC callback function type.
//
// For more details, see the MSDN.
//
// https://docs.microsoft.com/en-us/windows/win32/winmsg/using-hooks
type HOOKPROC func(code int32, wParam, lParam uintptr) uintptr
