package db

import "redis-clone/resp"

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
