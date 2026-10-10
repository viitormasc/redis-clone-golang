package main

import (
	"bufio"
	"fmt"
	"os"
	"redis-clone/handler"
	"redis-clone/resp"
	"strings"
)

func handleCommand(args handler.Arguments) string {
	cmd := strings.ToUpper(args[0])
	err := handler.CheckNumberOfArgs(args, cmd)
	if err != nil {
		return err.Error()
	}
	switch cmd {
	case "PING":
		return args.HandlePing()

	case "ECHO":
		return args.HandleEcho()

	case "COMMAND":
		return args.HandleCommand()
	case "SET":
		return args.HandleSet()
	case "GET":
		return args.HandleGet()
	}

	return fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd)
}

func encodeBulkString(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func readStdin() string {
	scanner := bufio.NewScanner(os.Stdin)
	var stream string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "*") {

		}
		stream += line + "\r\n"
	}
	return stream
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		stream := resp.ParseClientCommand(line)
		if stream == "" {
			continue
		}
		if strings.HasPrefix(stream, "*") {
			args := resp.DecodeBulkString(stream)
			response := handleCommand(args)
			fmt.Print(response)
			continue
		}
	}
}
