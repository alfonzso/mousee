package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	// "os/signal"

	"github.com/alfonzso/mousee/client"
	"github.com/alfonzso/mousee/common"
	"github.com/alfonzso/mousee/server"
	"github.com/moutend/go-hook/pkg/types"
	// "github.com/moutend/go-hook/pkg/types"
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
	cli, _ := Flags()
	// cli, update := Flags()

	// if update {
	// 	client.UpdateMode()
	// 	os.Exit(0)
	// } else if !cli {
	// 	go server.StartUpdateServer()
	// }

	// for {
	// 	time.Sleep(100 * time.Millisecond)
	// }

	if cli {
		client.WsClientMode()
	} else {
		if err := serverMode(); err != nil {
			log.Fatal(err)
		}
	}
}

func serverMode() error {

	mouseChan := make(chan types.MouseEvent, 100)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	infoLogger.Println("Server mode active ...")

	server := server.ServeWS()
	// server.ServeWS()

	// u := server.UdpConfig{
	// 	Addr: &net.UDPAddr{
	// 		Port: 1234,
	// 		IP:   net.ParseIP("192.168.1.100"),
	// 	},
	// 	ClientConnected: make(chan bool),
	// }

	// for {
	// 	time.Sleep(100 * time.Millisecond)
	// }

	go Mouse(nil, mouseChan)

	// // go MousePosHook(&u, signalChan, mouseChan)
	go MousePosHook(server, signalChan, mouseChan)

	// // server.StartServer()

	for {
		<-signalChan
		fmt.Println("Received shutdown signal")
		return nil
	}

	return nil
}
