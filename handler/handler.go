package handler

import (
	"fmt"
	"redis-clone/db"
	"redis-clone/resp"
	"strings"
)

type Arguments []string
type argsArity struct {
	min     int
	max     int
	Handler func(args Arguments) string
}

var CommandDef = map[string]argsArity{
	"PING": {
		min:     1,
		max:     2,
		Handler: HandlePing,
	},
	"ECHO": {
		min:     2,
		max:     2,
		Handler: HandleEcho,
	},
	"COMMAND": {
		min:     0,
		max:     2,
		Handler: HandleCommand,
	},
	"SET": {
		min:     3,
		max:     3,
		Handler: HandleSet,
	},
	"GET": {
		min:     2,
		max:     2,
		Handler: HandleGet,
	},
	"DBSIZE": {
		min:     1,
		max:     1,
		Handler: HandleDbSize,
	},
}

func emitCommandNotFound(cmd string) error {
	message := fmt.Sprintf("unknown command '%s'", cmd)
	errMessage := resp.Err(message)
	err := fmt.Errorf(errMessage)
	return err
}
func CheckNumberOfArgs(args Arguments, cmd string) error {
	numberOfArgs, ok := CommandDef[cmd]
	if !ok {
		return emitCommandNotFound(cmd)
	}
	message := fmt.Sprintf("wrong number of arguments for '%v' command", strings.ToUpper(cmd))
	if len(args) < numberOfArgs.min || len(args) > numberOfArgs.max {
		errMessage := resp.Err(message)
		err := fmt.Errorf(errMessage)
		return err
	}
	return nil
}

func HandlePing(args Arguments) string {
	if len(args) == 1 {
		return resp.SimpleString("PONG")
	} else {
		return resp.BulkString(args[1:])
	}
}

func HandleEcho(args Arguments) string {
	return resp.BulkString(args[1:])
}
func HandleCommand(args Arguments) string {
	return resp.SimpleString("OK")
}

func HandleSet(args Arguments) string {
	key, value := args[1], args[2]
	return db.SetValue(key, value)
}

func HandleGet(args Arguments) string {
	key := args[1]
	return db.GetValue(key)
}

func HandleDbSize(args Arguments) string {
	return db.Size()
}
