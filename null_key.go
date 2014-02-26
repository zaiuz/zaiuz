package zaiuz

import "code.google.com/p/go-uuid/uuid"
import "database/sql"

// Analog to sql.NullString for Key
type NullKey struct{ sql.NullString }

func NewNullKey() NullKey {
	return NullKey{sql.NullString{uuid.New(), true}}
}

func (nk NullKey) Key() Key {
	return Key(nk.String)
}
