package client

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"net"
	"os"
	"strings"

	"github.com/alfonzso/mousee/common"
)

// var reader *bufio.Reader
func readAsBytes(reader *bufio.Reader, args ...int) ([]byte, error) {
	size := 1024
	if len(args) > 0 {
		size = args[0]
	}
	p := make([]byte, size)
	readLen, err := reader.Read(p)
	if err != nil {
		return nil, err
	}
	return p[:readLen], nil
}
func readAsString(reader *bufio.Reader, args ...int) (string, error) {
	// size := 1024
	// if len(args) > 0 {
	// 	size = args[0]
	// }
	// p := make([]byte, size)
	data, err := readAsBytes(reader, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func readJson(data []byte, v any) {
	if err := json.Unmarshal(data, v); err != nil {
		panic(err)
	}
}

func removeLastNChar(s string, charCount int) string {
	return s[:len(s)-charCount]
}

func initFileProps(reader *bufio.Reader) (*os.File, common.UpdateData) {
	// var updateData common.UpdateData
	updateData := common.UpdateData{}
	message, err := readAsString(reader)
	if err != nil {
		panic(err)
	}
	readJson([]byte(message), &updateData)
	if updateData == (common.UpdateData{}) {
		panic(errors.New("UpdateData cannot be empty: " + message))
	}
	f := common.UpdateFile(updateData.AppVersion + "." + updateData.AppName)
	updateData.AppName = f.Name()
	return f, updateData
}

func getCrc32(updateData common.UpdateData) uint32 {
	dat, err := os.ReadFile(updateData.AppName)
	common.Check(err)

	crc32q := crc32.MakeTable(0xD5828281)
	return crc32.Checksum(dat, crc32q)
}

func UpdateMode() {

	infoLogger := log.New(os.Stdout, "INFO: ", 0)

	infoLogger.Println("Client mode active ...")
	// p := make([]byte, 1024)
	conn, err := net.Dial("tcp", "192.168.1.100:1235")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	reader := bufio.NewReader(conn)
	// readLen, err := reader.Read(p)
	// message := strings.TrimSpace(string(p[:readLen]))

	message, err := readAsString(reader)

	if message != "SUP" || err != nil {
		infoLogger.Println("EXITING, no SUP :'( ...")
	}

	fmt.Fprintf(conn, "UPDATE\n")
	message, err = readAsString(reader)

	for err == nil {
		if message == common.BeginUpdate() {
			infoLogger.Println("Begin of update")
			f, updateData := initFileProps(reader)

			d, err := readAsBytes(reader, 1024*1024)
			for err == nil {
				if strings.Contains(string(d), common.EndUpdate()) {
					finalBytes := removeLastNChar(string(d), len(common.EndUpdate()))
					f.Write([]byte(finalBytes))
					break
				}
				f.Write(d)
				d, err = readAsBytes(reader, 1024*1024)
			}
			f.Close()

			if crc := getCrc32(updateData); crc != updateData.AppCrc32 {
				panic(errors.New("CRC32 check failed... panicking"))
			}

			infoLogger.Println("End of update")
			break
		}

		message, err = readAsString(reader)
		infoLogger.Println(message)

	}

	conn.Close()
}
