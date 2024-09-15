package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/alfonzso/mousee/common"
)

// var reader *bufio.Reader
func readFrom(reader *bufio.Reader, args ...int) (string, error) {
	size := 1024
	if len(args) > 0 {
		size = args[0]
	}
	p := make([]byte, size)
	readLen, err := reader.Read(p)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(p[:readLen])), nil
}

func readJson(data []byte, v any) {
	if err := json.Unmarshal(data, v); err != nil {
		panic(err)
	}
}

func removeLastNChar(s string, charCount int) string {
	return s[:len(s)-4]
}

func initFileProps(reader *bufio.Reader) *os.File {
	var updateData common.UpdateData
	message, err := readFrom(reader)
	if err != nil {
		panic(err)
	}
	readJson([]byte(message), &updateData)
	f := common.UpdateFile(updateData.AppVersion + "." + updateData.FileName)
	return f
}

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
	// readLen, err := reader.Read(p)
	// message := strings.TrimSpace(string(p[:readLen]))

	message, err := readFrom(reader)

	if message != "SUP" || err != nil {
		infoLogger.Println("EXITING, no SUP :'( ...")
	}

	fmt.Fprintf(conn, "UPDATE\n")
	// readLen, err = reader.Read(p)
	message, err = readFrom(reader)

	for err == nil {
		// message := strings.TrimSpace(string(p[:readLen]))

		if message == common.BeginUpdate() {
			infoLogger.Println("Begin of update")
			// message, err = readFrom(reader)
			// readJson([]byte(message), &updateData)
			// f := common.UpdateFile(updateData.AppVersion + "." + updateData.FileName)
			f := initFileProps(reader)

			counter := 0
			for err == nil {
				// p = make([]byte, 1024*1024)
				// readLen, err = reader.Read(p)
				message, err = readFrom(reader)
				if strings.Contains(message, common.EndUpdate()) {
					// finalBytes := p[:readLen-len(common.EndUpdate())]
					finalBytes := removeLastNChar(message, len(common.EndUpdate()))
					f.Write([]byte(finalBytes))
					break
				}
				f.Write([]byte(message))
				counter += 1
			}
			f.Close()
			infoLogger.Println("End of update")
			break
		}

		message, err = readFrom(reader)
		// p = make([]byte, 1024)
		// readLen, err = reader.Read(p)
		infoLogger.Println(message)

	}

	conn.Close()
}
