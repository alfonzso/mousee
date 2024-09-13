package server

import (
	"bufio"
	"fmt"
	"net"
)

type TcpConfig struct {
	Addr            *net.UDPAddr
	Conn            *net.UDPConn
	Remoteaddr      *net.UDPAddr
	ClientConnected chan bool
}

// func (tcp *TcpConfig) ServeUDP() bool {
// 	ser, err := net.ListenUDP("udp", tcp.Addr)
// 	if err != nil {
// 		fmt.Printf("Some error %v\n", err)
// 		return false
// 	}
// 	tcp.Conn = ser
// 	return true
// }

func (tcp *TcpConfig) StartServer() {
	// Start a listener on port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 8080...")

	for {
		// Accept a connection from the client
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the client in a separate goroutine
		go tcp.handleConnection(conn)
	}
}

func (tcp *TcpConfig) handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	// Read data from the client
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			return
		}

		// Print the message received from the client
		fmt.Printf("Received: %s", message)

		// Send a response back to the client
		conn.Write([]byte("Message received\n"))
	}
}
