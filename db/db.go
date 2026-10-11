package db

import (
	"redis-clone/resp"
	"strconv"
)

type KvPair map[string]string

var KeyValuePair = KvPair{}

func SetValue(k string, v string) string {
	KeyValuePair[k] = v
	return resp.SimpleString("OK")
}

func GetValue(k string) string {
	value, ok := KeyValuePair[k]
	if ok {
		arrStr := []string{value}
		return resp.BulkString(arrStr)
	}
	return resp.BulkString(nil)
}

func Size() string {
	return resp.Int(len(KeyValuePair))
}

func ChangeInt(k string, incr bool, difValue string) string {
	value, ok := KeyValuePair[k]
	if !ok {
		value = "0"
	}
	valueInt, err := strconv.Atoi(value)
	diffValueInt, errInc := strconv.Atoi(difValue)
	if err != nil || errInc != nil {
		return resp.Err("value is not an integer or out of range")
	}
	if incr {
		valueInt += diffValueInt
	} else {
		valueInt -= diffValueInt
	}
	KeyValuePair[k] = strconv.Itoa(valueInt)
	return resp.Int(valueInt)
}
