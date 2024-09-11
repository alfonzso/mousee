package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

	_, err = reader.Read(p)

	for err != io.EOF {
		fmt.Printf("%s\n", p)
		p = make([]byte, 1024)
		_, err = reader.Read(p)
	}

	conn.Close()
}
