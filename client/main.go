package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
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

	// if err == nil {
	// 	fmt.Printf("%s\n", p)
	// 	fmt.Printf("%d\n", kek)
	// 	// kek, err = reader.Read(p)
	// 	// fmt.Printf("%s\n", p)
	// 	// fmt.Printf("%d\n", kek)
	// 	// kek, err = reader.Read(p)
	// 	// fmt.Printf("%s\n", p)
	// 	// fmt.Printf("%d\n", kek)
	// 	// kek, err = reader.Read(p)
	// 	// fmt.Printf("%s\n", p)
	// 	// fmt.Printf("%d\n", kek)
	// 	// kek, err = bufio.NewReader(conn).Read(p)
	// 	// fmt.Printf("%s\n", p)
	// 	// fmt.Printf("%d\n", kek)
	// 	// kek, err = bufio.NewReader(conn).Read(p)
	// 	// fmt.Printf("%s\n", p)
	// 	// fmt.Printf("%d\n", kek)
	// } else {
	// 	fmt.Printf("Some error %v\n", err)
	// }
	conn.Close()
}
