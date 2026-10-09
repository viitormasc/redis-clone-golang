package resp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func readLine(buf *bytes.Buffer) string {
	var readLine strings.Builder
	for {
		b, err := buf.ReadByte()
		readStr := readLine.String()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println("Error reading byte:", err)
			break
		}

		if b == '\n' && strings.HasSuffix(readStr, "\r") {
			return readStr[:len(readStr)-1]
		}
		readLine.WriteByte(b)
	}
	return readLine.String()
}

func DecodeBulkString(stream string) []string {
	buf := bytes.NewBufferString(stream)
	header := readLine(buf)
	numberOfLines, err := strconv.Atoi(header[1:])
	if err != nil {
		fmt.Println("Error converting string to int", err)
	}
	args := []string{}

	for range numberOfLines {
		bulkReader := readLine(buf)
		// gets the length of the next argument
		length, err := strconv.Atoi(string(bulkReader[1:]))
		if err != nil {
			fmt.Println("Error converting string to int", err)
		}
		//read exactly the byte size of the argument to avoid getting junk
		dataBuf := make([]byte, length)
		_, err = io.ReadFull(buf, dataBuf)
		data := string(dataBuf)
		if err != nil {
			fmt.Println("Error reading byte:", err)
		}
		//get rid of the next /r/n which is 2 bytes long
		for range 2 {
			buf.ReadByte()
		}
		args = append(args, string(data))
	}
	return args
}
