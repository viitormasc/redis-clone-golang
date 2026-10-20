package handler

import (
	"fmt"
	"redis-clone/db"
	"redis-clone/resp"
	"slices"

	// "slices"
	"strconv"
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
	"EXPIRE": {
		min:     3,
		max:     3,
		Handler: HandleExpire,
	},
	"TTL": {
		min:     2,
		max:     2,
		Handler: HandleTtl,
	},
	"PTTL": {
		min:     2,
		max:     2,
		Handler: HandlePttl,
	},
	"PERSIST": {
		min:     2,
		max:     2,
		Handler: HandlePersist,
	},
	"INCR": {
		min:     2,
		max:     2,
		Handler: HandleIncr,
	},
	"INCRBY": {
		min:     3,
		max:     3,
		Handler: HandleIncrby,
	},
	"DECR": {
		min:     2,
		max:     2,
		Handler: HandleDecr,
	},
	"DECRBY": {
		min:     3,
		max:     3,
		Handler: HandleDecrby,
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
		max:     6,
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
	keyExists := db.GetValue(key) != resp.NullBulk()
	var ttlSec int
	if len(args) >= 3 {
		flags := args[3:]
		for i, flag := range flags {
			flags[i] = strings.ToUpper(flag)
		}
		if slices.Contains(flags, "EX") && slices.Contains(flags, "PX") {
			return resp.Err("EX and PX arguments can't be used together")
		}
		for i, flag := range flags {
			upperFlag := flag
			if upperFlag == "NX" && keyExists {
				return resp.NullBulk()
			}
			if upperFlag == "XX" && !keyExists {
				return resp.NullBulk()
			}
			if upperFlag == "EX" {
				time, err := strconv.Atoi(flags[i+1])
				if err != nil {
					return resp.Err("Time should be a number")
				}
				ttlSec = time
			}
			if upperFlag == "PX" {
				time, err := strconv.Atoi(flags[i+1])
				if err != nil {
					return resp.Err("Time should be a number")
				}
				ttlSec = time / 1000
			}
		}
	}
	response := db.SetValue(key, value)
	if ttlSec != 0 {
		db.ChangeExpireTime(key, strconv.Itoa(ttlSec))
	}
	return response
}

func HandleGet(args Arguments) string {
	key := args[1]
	return db.GetValue(key)
}

func HandleDbSize(args Arguments) string {
	return db.Size()
}

func HandleIncr(args Arguments) string {
	key := args[1]
	return db.ChangeInt(key, true, "1")
}

func HandleIncrby(args Arguments) string {
	key, incrValue := args[1], args[2]
	return db.ChangeInt(key, true, incrValue)
}

func HandleDecr(args Arguments) string {
	key := args[1]
	return db.ChangeInt(key, false, "1")
}

func HandleExpire(args Arguments) string {
	key, ttlSec := args[1], args[2]
	return db.ChangeExpireTime(key, ttlSec)
}
func HandleTtl(args Arguments) string {
	key := args[1]
	return db.GetExpireTime(key, db.Seconds)
}
func HandlePttl(args Arguments) string {
	key := args[1]
	return db.GetExpireTime(key, db.Milliseconds)
}
func HandlePersist(args Arguments) string {
	key := args[1]
	return db.DeleteExpireKey(key)
}

func HandleDecrby(args Arguments) string {
	key, incrValue := args[1], args[2]
	return db.ChangeInt(key, false, incrValue)
}
