package db

type Pair struct {
	key   string
	value any
}

type DB []Pair

var Database = DB{}

func (db *DB) SetValue(k string, v string) {

	kv := Pair{key: k, value: v}
	Database = append(Database, kv)
}

func (db *DB) GetValue(k string) any {
	for _, kv := range Database {
		if kv.key == k {
			return kv.key
		}
	}
	return nil
}
