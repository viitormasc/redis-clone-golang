package handler

import (
	"fmt"
	"redis-clone/db"
	"redis-clone/resp"
	"strings"
)

type argsArity struct {
	min int
	max int
}

var argsCount = map[string]argsArity{
	"PING": {
		min: 1,
		max: 2,
	},
	"ECHO": {
		min: 2,
		max: 2,
	},
	"COMMAND": {
		min: 0,
		max: 2,
	},
	"SET": {
		min: 3,
		max: 3,
	},
	"GET": {
		min: 2,
		max: 2,
	},
}

func emitCommandNotFound(cmd string) error {
	message := fmt.Sprintf("unknown command '%s'", cmd)
	errMessage := resp.Err(message)
	err := fmt.Errorf(errMessage)
	return err
}
func CheckNumberOfArgs(args Arguments, cmd string) error {
	numberOfArgs, ok := argsCount[cmd]
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

type Arguments []string

func (args Arguments) HandlePing() string {
	if len(args) == 1 {
		return resp.SimpleString("PONG")
	} else {
		return resp.BulkString(args[1:])
	}
}

func (args Arguments) HandleEcho() string {
	return resp.BulkString(args[1:])
}
func (args Arguments) HandleCommand() string {
	return resp.SimpleString("OK")
}

func (args Arguments) HandleSet() string {
	key, value := args[1], args[2]
	return db.SetValue(key, value)
}

func (args Arguments) HandleGet() string {
	key := args[1]
	return db.GetValue(key)
}
