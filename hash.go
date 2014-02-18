package zaiuz

import "code.google.com/p/go.crypto/bcrypt"
import "fmt"
import "database/sql/driver"

// Represents a hashed string usually used for storing passwords. This type implements
// driver.Valuer and sql.Scanner so you can use this as your model's password hash type
// directly. The underlying algorithm use google's bcrypt library.
type Hash string

// Returns a new hash from a string.
func NewHash(original string) Hash {
	bytes := []byte(original)
	bytes, e := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if e != nil {
		panic(e) // TODO: Hmm? Not sure if should panic here.
	}

	return Hash(string(bytes))
}

// Returns wether or not the receiver hash is the hash of the given original string. Use
// this method to compare user's password input with stored hash value.
func (hash Hash) MatchOriginal(original string) bool {
	h, o := []byte(hash), []byte(original)
	e := bcrypt.CompareHashAndPassword(h, o)
	return e == nil
}

func (hash *Hash) Scan(src interface{}) (e error) {
	switch src.(type) {
	case string:
		*hash = NewHash(src.(string))
	case []byte:
		*hash = NewHash(string(src.([]byte)))
	default:
		e = fmt.Errorf("cannot convert driver value `%v` to a Hash value.", src)
	}
	return
}

func (hash Hash) Value() (driver.Value, error) {
	return driver.Value(string(hash)), nil
}
