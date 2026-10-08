package resp

import (
	"fmt"
	"strconv"
)

const crlf string = "\r\n"
const marker string = "$"

func SimpleString(s string) string {
	simpleResp := "+" + s + crlf
	return simpleResp
}

func BulkString(str []string) string {
	bulkResp := ""
	for _, s := range str {
		sLen := len(s)
		lenString := strconv.Itoa(sLen)
		bulkResp = marker + lenString + crlf + s + crlf
	}
	return bulkResp
}

func Err(message string) string {
	errResp := fmt.Sprintf("-ERR %s\r\n", message)
	return errResp
}

func NullBulk(s string) string {
	return "$-1\r\n"
}

func Int(int int64) string {
	errResp := fmt.Sprintf(":%v\r\n", int)
	return errResp
}
