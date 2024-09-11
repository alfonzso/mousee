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
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	// "time"

	// "os"
	// "os/signal"
	// "time"

	"github.com/alfonzso/mousee/client"
	"github.com/alfonzso/mousee/server"
)

func Flags() bool {

	var client bool

	flag.BoolVar(&client, "client", false, "Client or Server mode, default Server")

	flag.Parse()

	return client
}

var infoLogger = log.New(os.Stdout, "INFO: ", 0)

func main() {
	// log.SetFlags(0)
	// log.SetPrefix("error: ")
	// infoLogger.Println("Client mode active ...")

	if cli := Flags(); cli {
		client.ClientMode()
	} else {
		if err := serverMode(); err != nil {
			log.Fatal(err)
		}
	}
}

func serverMode() error {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	infoLogger.Println("Server mode active ...")

	u := server.UdpConfig{
		Addr: &net.UDPAddr{
			Port: 1234,
			IP:   net.ParseIP("192.168.1.100"),
		},
		ClientConnected: make(chan bool),
	}

	go Mouse()

	go MousePosHook(&u, signalChan)

	u.StartServer()

	for {
		// serve forever
		// time.Sleep(100 * time.Millisecond)
		// select {
		// case
		<-signalChan
		fmt.Println("Received shutdown signal")
		return nil
		// }
	}

	return nil
}
