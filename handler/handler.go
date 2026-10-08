package handler

import "redis-clone/resp"

type Arguments []string

func (args Arguments) HandlePing() string {
	if len(args) == 1 {
		return resp.SimpleString("PONG")
	} else {
		return resp.BulkString(args[1:])
	}
}
