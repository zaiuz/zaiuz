package zaiuz

import "code.google.com/p/go-uuid/uuid"
import "fmt"
import "database/sql/driver"

// Represents a GUID database key. This type implements driver.Valuer and sql.Scanner so
// you can use this as your model's Key type directly.
type Key string

// Returns a new GUID key.
func NewKey() Key {
	return Key(uuid.New())
}

func (key *Key) Scan(src interface{}) (e error) {
	switch src.(type) {
	case string:
		*key = Key(src.(string))
	case []byte:
		*key = Key(string(src.([]byte)))
	default:
		e = fmt.Errorf("cannot convert driver value `%v` to a Key value.", src)
	}
	return
}

// TODO: validate
func (key Key) Value() (driver.Value, error) {
	return driver.Value(string(key)), nil
}
