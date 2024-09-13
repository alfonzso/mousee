package server

import (
	"fmt"
	"net"
	"os"
	"time"
)

type UdpConfig struct {
	Addr            *net.UDPAddr
	Conn            *net.UDPConn
	Remoteaddr      *net.UDPAddr
	ClientConnected chan bool
}

// func (u UdpConfig) SendResponse(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
func (u UdpConfig) SendResponse(msg string) {
	_, err := u.Conn.WriteToUDP([]byte(msg), u.Remoteaddr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
	// fmt.Printf("%d %d \r", len, err)
}

func (u *UdpConfig) ServeUDP() bool {
	ser, err := net.ListenUDP("udp", u.Addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return false
	}
	u.Conn = ser
	return true
}

// The function will block the runtime until someone connects to us.
func (u UdpConfig) IsClientConnected(signalChan chan os.Signal) bool {
	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case connected := <-u.ClientConnected:
			if connected {
				fmt.Println(">> ClientConnected")
				return true
			}

		case <-signalChan:
			fmt.Println(">> App shutting down")
			return false
		}
	}
}

func (u *UdpConfig) StartServer() {
	p := make([]byte, 2048)

	u.ServeUDP()
	u.ClientConnected <- false
	readLen, remoteaddr, err := u.Conn.ReadFromUDP(p)

	fmt.Printf("Read a message from %v %s \n", remoteaddr, p[:readLen])
	if err != nil {
		fmt.Printf("Some error  %v", err)
	}

	u.Remoteaddr = remoteaddr

	fmt.Println("Starting ...", u.Remoteaddr)
	u.ClientConnected <- true

}
