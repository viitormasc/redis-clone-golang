package resp

import (
	"fmt"
	"strconv"
	"strings"
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

func Int(int int) string {
	errResp := fmt.Sprintf(":%v\r\n", int)
	return errResp
}

func parseArgs(line string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	for _, ch := range line {
		switch {
		case ch == '"' && !inQuotes:
			inQuotes = true
		case ch == '"' && inQuotes:
			if current.String() == "" {
				args = append(args, "")
			}
			inQuotes = false
		case ch == ' ' && !inQuotes:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
func ParseClientCommand(line string) string {
	args := parseArgs(line)
	argsNum := len(args)
	encodeResp := fmt.Sprintf("*%v\r\n", argsNum)
	for _, arg := range args {
		argLen := len(arg)
		encodeResp += fmt.Sprintf("$%v\r\n%v\r\n", argLen, arg)
	}
	return encodeResp
}
