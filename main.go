// //go:build windows

// package main

// import (
//     "fmt"
//     "syscall"
//     "unsafe"
// )

// func main() {
//     userDll := syscall.NewLazyDLL("user32.dll")
//     getWindowRectProc := userDll.NewProc("GetCursorPos")
//     type POINT struct {
//         X, Y int32
//     }
//     var pt POINT
//     _, _, eno := syscall.SyscallN(getWindowRectProc.Addr(), uintptr(unsafe.Pointer(&pt)))
//     if eno != 0 {
//         fmt.Println(eno)
//     }
//     fmt.Printf("[cursor.Pos] X:%d Y:%d\n", pt.X, pt.Y)
// }
//go:build windows
// +build windows

package main

import (
	// "fmt"
	"log"
	"net"

	// "time"

	// "os"
	// "os/signal"
	// "time"

	"github.com/alfonzso/mousee/server"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	u := server.UdpConfig{
		Addr: &net.UDPAddr{
			Port: 1234,
			IP:   net.ParseIP("192.168.1.100"),
		},
		ClientConnected: make(chan bool),
	}

	// go Mouse()

	go MousePosHook(u)

	u.StartServer()
	// for {
	// 	// serve forever
	// 	time.Sleep(1 * time.Millisecond)
	// }

	return nil
}
