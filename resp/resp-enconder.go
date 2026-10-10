package resp

import (
	"fmt"
	"strconv"
)

const crlf string = "\r\n"
const marker string = "$"

func SimpleString(s string) string {
	simpleResp := fmt.Sprintf("+%v%v", s, crlf)
	return simpleResp
}

func BulkString(arrStr []string) string {
	if len(arrStr) == 0 {
		return NullBulk()
	}
	bulkResp := ""
	for _, s := range arrStr {
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

func NullBulk() string {
	return "$-1\r\n"
}

func Int(int int64) string {
	errResp := fmt.Sprintf(":%v\r\n", int)
	return errResp
}
