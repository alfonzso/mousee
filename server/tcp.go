package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/alfonzso/mousee/common"
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

		// Print the message received from the client
		fmt.Printf("Received: %s", message)

		if message == "UPDATE" {
			message = ""
			// fmt.Println("Goootttt ittttt")
			conn.Write([]byte(common.BeginUpdate()))
			dat, err := os.ReadFile("mousee.exe")
			common.Check(err)
			fmt.Printf(">>>>>>>> %d\n", len(dat))
			n, e := conn.Write(dat)
			// n, e := conn.Write(dat+ []byte("END_UPDATE"))
			// n, e := conn.Write(append(dat, []byte("END_UPDATE")...))
			fmt.Printf(">>>>>>>> %d %v\n", n, e)
			conn.Write([]byte(common.EndUpdate()))
		}

		// Send a response back to the client
		// conn.Write([]byte("Message received\n"))
	}
}
