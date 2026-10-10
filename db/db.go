package db

type Pair struct {
	key   string
	value string
}

type DB []Pair

var Database = DB{}

func (db *DB) SetValue(k string, v string) {
	// vArr := db.GetValue(k)
	// vArr = append(vArr, v)
	kv := Pair{key: k, value: v}
	Database = append(Database, kv)
}

func (db *DB) GetValue(k string) string {
	for _, kv := range Database {
		if kv.key == k {
			return kv.value
		}
	}
	return ""
}
