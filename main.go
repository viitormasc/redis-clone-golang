package main

import (
	"bufio"
	"fmt"
	// "io"
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
	// stream := readStdin()
	scanner := bufio.NewScanner(os.Stdin)
	var stream string
	run := 1
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "*") && run != 1 {
			args := resp.DecodeBulkString(stream)
			response := handleCommand(args)
			fmt.Print(response)
			stream = line + "\r\n"
			continue
		}
		stream += line + "\r\n"
		run++
	}
	args := resp.DecodeBulkString(stream)
	response := handleCommand(args)
	fmt.Print(response)
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
