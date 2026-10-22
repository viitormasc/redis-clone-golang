package db

import (
	"redis-clone/resp"
	"strconv"
	"time"
)

type KvPair map[string]string

var keyValuePair = KvPair{}
var expires = KvPair{}

type Unit string

const (
	Milliseconds Unit = "ms"
	Seconds      Unit = "sec"
)

func SetValue(k string, v string) string {
	keyValuePair[k] = v
	return resp.SimpleString("OK")
}

func GetValue(k string) string {
	deleteIfExpired(k)
	value, ok := keyValuePair[k]
	if ok {
		arrStr := []string{value}
		return resp.BulkString(arrStr)
	}
	return resp.BulkString(nil)
}

func Size() string {
	size := 0
	for k, _ := range keyValuePair {
		if !isExpired(k) {
			size++
		}
	}
	return resp.Int(size)
}

func ChangeInt(k string, incr bool, difValue string) string {
	deleteIfExpired(k)
	value, ok := keyValuePair[k]
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
	keyValuePair[k] = strconv.Itoa(valueInt)
	return resp.Int(valueInt)
}

func DeleteKvKey(k string) string {
	_, ok := keyValuePair[k]
	if !ok {
		return resp.Int(0)
	}
	delete(keyValuePair, k)
	return resp.Int(1)
}

func DeleteExpireKey(k string) string {
	_, ok := expires[k]
	if !ok {
		return resp.Int(0)
	}
	delete(expires, k)
	return resp.Int(1)
}

func ChangeExpireTime(k string, ttlSecStr string) string {
	addTime, err := strconv.Atoi(ttlSecStr)
	if err != nil {
		return resp.Err("time must be an Int in seconds")
	}
	_, ok := keyValuePair[k]
	if addTime == 0 || !ok {
		return DeleteKvKey(k)
	}
	now := now()

	expiresAt := now + addTime

	expires[k] = strconv.Itoa(expiresAt)
	return resp.Int(1)
}

func now() int {
	return int(time.Now().UTC().UnixMilli() / 1000)
}
func isExpired(k string) bool {
	ttlStr, ok := expires[k]
	if !ok {
		return true
	}
	ttl, _ := strconv.Atoi(ttlStr)
	now := now()

	return ttl < now

}
func deleteIfExpired(k string) {
	expired := isExpired(k)
	if expired {
		delete(keyValuePair, k)
		delete(expires, k)
	}
}

func Wait(waitTimeMs string) string {
	waitTimeMsInt, err := strconv.Atoi(waitTimeMs)
	if err != nil {
		return resp.Err("time must be a integer")
	}
	time.Sleep(time.Duration(waitTimeMsInt) * time.Millisecond)
	return resp.SimpleString("OK")
}

func GetExpireTime(k string, unit Unit) string {
	if unit != Milliseconds && unit != Seconds {
		return resp.Err("you must specify the unity to get the answer")
	}
	_, ok := keyValuePair[k]
	if !ok {
		return resp.Int(-2)
	}

	expireTime, ok := expires[k]
	if !ok {
		return resp.Int(-1)
	}
	expireTInt, err := strconv.Atoi(expireTime)
	if err != nil {
		return resp.Err("time must be an Int in seconds")
	}
	remaining := expireTInt - now()
	if unit == Milliseconds {
		remaining = remaining * 1000
	}
	if remaining >= 0 {
		return resp.Int(remaining)
	} else {
		return resp.Int(-2)
	}
}

func Exists(k string) string {
	deleteIfExpired(k)
	_, ok := keyValuePair[k]
	if ok {
		return resp.Int(1)
	}
	return resp.Int(0)
}
