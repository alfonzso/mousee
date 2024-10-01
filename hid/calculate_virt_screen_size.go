package hid

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/alfonzso/mousee/common"
)

var (
	procEnumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	procGetMonitorInfo      = user32.NewProc("GetMonitorInfoW")
)

func CalculateScreen() common.VScreenSize {
	// Variables to calculate the overall screen dimensions
	minLeft, minTop := int32(0), int32(0)
	maxRight, maxBottom := int32(0), int32(0)
	isFirstMonitor := true

	// Call EnumDisplayMonitors to get information about all monitors
	_, _, _ = procEnumDisplayMonitors.Call(
		0,
		0,
		syscall.NewCallback(func(hMonitor, hdcMonitor, lprcMonitor, dwData uintptr) uintptr {
			// Fill MonitorInfo structure
			var mi common.MonitorInfo
			mi.Size = uint32(unsafe.Sizeof(mi))

			ret, _, _ := procGetMonitorInfo.Call(hMonitor, uintptr(unsafe.Pointer(&mi)))
			if ret != 0 {
				// Get monitor position
				left := mi.Monitor.Left
				top := mi.Monitor.Top
				right := mi.Monitor.Right
				bottom := mi.Monitor.Bottom

				fmt.Printf("Monitor: %dx%d at position (Left: %d, Top: %d, Right: %d, Bottom: %d)\n",
					right-left, bottom-top, left, top, right, bottom)

				// Initialize with the first monitor's values
				if isFirstMonitor {
					minLeft, minTop = left, top
					maxRight, maxBottom = right, bottom
					isFirstMonitor = false
				} else {
					// Update the minimum and maximum coordinates
					if left < minLeft {
						minLeft = left
					}
					if top < minTop {
						minTop = top
					}
					if right > maxRight {
						maxRight = right
					}
					if bottom > maxBottom {
						maxBottom = bottom
					}
				}
			}
			return 1 // Continue enumeration
		}),
		0,
	)

	// Calculate overall virtual screen dimensions
	return common.VScreenSize{maxBottom - minTop, maxRight - minLeft}
}
