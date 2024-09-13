package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/alfonzso/mousee/common"

	"github.com/go-vgo/robotgo"
	_ "github.com/go-vgo/robotgo/base"
	_ "github.com/go-vgo/robotgo/key"
	_ "github.com/go-vgo/robotgo/mouse"
	_ "github.com/go-vgo/robotgo/screen"
	_ "github.com/go-vgo/robotgo/window"
)

func ClientMode() {
	infoLogger := log.New(os.Stdout, "INFO: ", 0)

	infoLogger.Println("Client mode active ...")
	p := make([]byte, 1024)
	conn, err := net.Dial("udp", "192.168.1.100:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")
	reader := bufio.NewReader(conn)

	readLen, err := reader.Read(p)

	var mouseData common.MouseData

	for err != io.EOF {
		// pOK := p[:readLen]
		if err := json.Unmarshal(p[:readLen], &mouseData); err != nil {
			panic(err)
		}
		// fmt.Println("mouseData", mouseData)
		robotgo.Move(int(mouseData.X), int(mouseData.Y))

		if uintptr(common.WM_LBUTTONDOWN) == mouseData.Msg {
			robotgo.Click("left")
		}

		if uintptr(common.WM_RBUTTONDOWN) == mouseData.Msg {
			robotgo.Click("right")
		}

		p = make([]byte, 1024)
		readLen, err = reader.Read(p)
	}

	conn.Close()
}
