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
	"time"

	// "time"

	// "os"
	// "os/signal"
	// "time"

	"github.com/alfonzso/mousee/client"
	"github.com/alfonzso/mousee/common"
	"github.com/alfonzso/mousee/server"
	"github.com/moutend/go-hook/pkg/types"
)

func Flags() (bool, bool) {

	var client bool
	var update bool
	var version bool

	flag.BoolVar(&client, "client", false, "Client or Server mode, default Server")
	flag.BoolVar(&update, "update", false, "Client will update itself from server")
	flag.BoolVar(&version, "version", false, "App version")

	flag.Parse()

	if version {
		fmt.Println(common.AppVersion)
		os.Exit(0)
	}

	return client, client && update
}

var infoLogger = log.New(os.Stdout, "INFO: ", 0)

func main() {
	_, update := Flags()
	// log.SetFlags(0)
	// log.SetPrefix("error: ")
	// infoLogger.Println("Client mode active ...")

	if update {
		client.UpdateMode()
		os.Exit(0)
	} else {
		go server.StartUpdateServer()
	}

	for {
		// serve forever
		time.Sleep(100 * time.Millisecond)
		// select {
		// case
		// }
	}

	// if cli {
	// 	client.ClientMode()
	// } else {
	// 	if err := serverMode(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
}

func serverMode() error {

	mouseChan := make(chan types.MouseEvent, 100)
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

	go Mouse(nil, mouseChan)

	go MousePosHook(&u, signalChan, mouseChan)

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
