package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"net"
	"os"
	"strings"

	"github.com/alfonzso/mousee/common"
	// "../common"
	// "os"
	// "github.com/alfonzso/mousee/common"
)

// type TcpConfig struct {
// 	Addr            *net.UDPAddr
// 	Conn            *net.UDPConn
// 	Remoteaddr      *net.UDPAddr
// 	ClientConnected chan bool
// }

// func (tcp *TcpConfig) ServeUDP() bool {
// 	ser, err := net.ListenUDP("udp", tcp.Addr)
// 	if err != nil {
// 		fmt.Printf("Some error %v\n", err)
// 		return false
// 	}
// 	tcp.Conn = ser
// 	return true
// }

// func (tcp *TcpConfig) StartUpdateServer() {
func StartUpdateServer() {
	// Start a listener on port 8080
	ln, err := net.Listen("tcp", ":1235")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 1235...")

	for {
		// Accept a connection from the client
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the client in a separate goroutine
		go handleConnection(conn)
	}
}

// func (tcp *TcpConfig) handleConnection(conn net.Conn) {
func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	conn.Write([]byte("SUP\n"))

	// Read data from the client
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			return
		}

		fmt.Printf("Received: %s \n", message)

		if message == "UPDATE" {
			message = ""
			conn.Write([]byte(common.BeginUpdate()))
			dat, err := os.ReadFile(common.AppName)
			common.Check(err)

			crc32q := crc32.MakeTable(0xD5828281)
			b, err := json.Marshal(
				common.UpdateData{
					AppName:    common.AppName,
					AppVersion: common.AppVersion,
					AppCrc32:   crc32.Checksum(dat, crc32q),
				},
			)
			if err == nil {
				conn.Write(b)
			}
			conn.Write(dat)
			conn.Write([]byte(common.EndUpdate()))
		}

		// Send a response back to the client
		// conn.Write([]byte("Message received\n"))
	}
}
