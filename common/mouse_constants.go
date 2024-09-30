package common

import "github.com/moutend/go-hook/pkg/types"

const (
	WH_JOURNALRECORD   types.Hook = 0
	WH_JOURNALPLAYBACK types.Hook = 1
	WH_KEYBOARD        types.Hook = 2
	WH_GETMESSAGE      types.Hook = 3
	WH_CALLWNDPROC     types.Hook = 4
	WH_CBT             types.Hook = 5
	WH_SYSMSGFILTER    types.Hook = 6
	WH_MOUSE           types.Hook = 7
	WH_DEBUG           types.Hook = 9
	WH_SHELL           types.Hook = 10
	WH_FOREGROUNDIDLE  types.Hook = 11
	WH_CALLWNDPROCRET  types.Hook = 12
	WH_KEYBOARD_LL     types.Hook = 13
	WH_MOUSE_LL        types.Hook = 14
)

const (
	WM_MOUSWHEELUP   types.Message = 0x0198
	WM_MOUSWHEELDOWN types.Message = 0x0199
	WM_MOUSEMOVE     types.Message = 0x0200
	WM_LBUTTONDOWN   types.Message = 0x0201
	WM_LBUTTONUP     types.Message = 0x0202
	WM_MBUTTONDOWN   types.Message = 0x0207
	WM_MBUTTONUP     types.Message = 0x0208
	WM_MOUSEWHEEL    types.Message = 0x020A
	WM_MOUSEHWHEEL   types.Message = 0x020E
	WM_RBUTTONDOWN   types.Message = 0x0204
	WM_RBUTTONUP     types.Message = 0x0205
	WM_KEYDOWN       types.Message = 0x0100
	WM_KEYUP         types.Message = 0x0101
	WM_SYSKEYDOWN    types.Message = 0x0104
	WM_SYSKEYUP      types.Message = 0x0105
)
