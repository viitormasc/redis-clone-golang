package db

import "redis-clone/resp"

type Pair struct {
	key   string
	value string
}

type DB []Pair

var Database = DB{}

func (db *DB) SetValue(k string, v string) string {
	// vArr := db.GetValue(k)
	// vArr = append(vArr, v)
	kv := Pair{key: k, value: v}
	Database = append(Database, kv)
	return resp.SimpleString("OK")
}

func (db *DB) GetValue(k string) string {
	for _, kv := range Database {
		if kv.key == k {
			arrStr := []string{kv.value}
			return resp.BulkString(arrStr)
		} else {
			return resp.NullBulk()
		}
	}
	return resp.NullBulk()
}
