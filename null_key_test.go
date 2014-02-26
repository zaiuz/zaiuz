package zaiuz_test

import "testing"
import "database/sql"
import "database/sql/driver"
import . "github.com/zaiuz/zaiuz"
import a "github.com/stretchr/testify/assert"

var nullKey NullKey = NewNullKey()
var _ sql.Scanner = &nullKey
var _ driver.Valuer = &nullKey
var _ driver.Valuer = nullKey

func TestNewNullKey(t *testing.T) {
	result := NewNullKey()
	empty := NullKey{}

	a.NotEqual(t, result, empty, "new null key returns empty key.")
}

func TestNullKey_Key(t *testing.T) {
	nk := NewNullKey()
	k := Key(nk.String)

	a.Equal(t, nk.Key(), k, "key method returns wrong key.")
}

func TestNullKey_Scan(t *testing.T) {
	var nullKey NullKey

	test := func(src interface{}, valid bool, expected string) {
		nullKey = NullKey{}
		e := nullKey.Scan(src)
		a.NoError(t, e)
		a.Equal(t, nullKey.Valid, valid,
			"scanned null key has wrong validity for value: `%v`.", src)
		if valid {
			a.Equal(t, nullKey.String, expected,
				"scanned null key has wrong value for: `%v`.", src)
		}
	}

	test(nil, false, "")
	test("string", true, "string")
	test([]byte("string"), true, "string")
}

func TestNullKey_Value(t *testing.T) {
	nk := NewNullKey()
	v, e := nk.Value()
	a.NoError(t, e)
	a.Equal(t, v, driver.Value(string(nk.String)), "wrong driver value for valid null key.")

	nk = NullKey{}
	v, e = nk.Value()
	a.NoError(t, e)
	a.Equal(t, v, driver.Value(nil), "driver value not nil for nil null key.")
}
