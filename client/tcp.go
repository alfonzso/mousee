package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/alfonzso/mousee/common"
)

func UpdateMode() {
	infoLogger := log.New(os.Stdout, "INFO: ", 0)

	infoLogger.Println("Client mode active ...")
	p := make([]byte, 1024)
	conn, err := net.Dial("tcp", "192.168.1.100:1235")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	reader := bufio.NewReader(conn)
	readLen, err := reader.Read(p)
	message := strings.TrimSpace(string(p[:readLen]))

	if message != "SUP" || err != nil {
		infoLogger.Println("EXITING, no SUP :'( ...")
	}

	fmt.Fprintf(conn, "UPDATE\n")

	readLen, err = reader.Read(p)

	// infoLogger.Println(string(p[:readLen]))

	// var mouseData common.MouseData

	// var

	for err == nil {
		// pOK := p[:readLen]
		// if err := json.Unmarshal(p[:readLen], &mouseData); err != nil {
		// 	panic(err)
		// }

		message := strings.TrimSpace(string(p[:readLen]))

		// infoLogger.Println("================> ", message)

		if message == common.BeginUpdate() {
			infoLogger.Println("Begin of update")
			f, err := os.OpenFile("fafa.exe", os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			// f, err := os.OpenFile("fafa.exe", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}

			// f.Truncate(0)
			// f.Seek(0, 0)

			// for message == "END_UPDATE" {
			counter := 0
			for err == nil {
				p = make([]byte, 1024*1024)
				readLen, err = reader.Read(p)
				// if err != nil {
				// 	panic(err)
				// }
				// message = strings.TrimSpace(string(p[:readLen]))
				// infoLogger.Println("eeeeeee", readLen)
				// fmt.Printf("> %d  %v     \n", counter, err)
				if strings.Contains(string(p[:readLen]), common.EndUpdate()) {
					// fmt.Printf("> %d  %v  %d   \n", counter, err, readLen)
					// fmt.Printf(">    %v     \n", string(p[:readLen]))
					finalBytes := p[:readLen-len(common.EndUpdate())]
					f.Write(finalBytes)
					break
				}
				f.Write(p[:readLen])
				counter += 1
			}
			f.Close()
			infoLogger.Println("End of update")
			break
		}

		// if err != nil {
		// 	break
		// }

		p = make([]byte, 1024)
		readLen, err = reader.Read(p)
		infoLogger.Println(string(p[:readLen]))

	}

	conn.Close()
}
