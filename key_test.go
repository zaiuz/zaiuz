package zaiuz

import "testing"
import "database/sql"
import "database/sql/driver"
import a "github.com/stretchr/testify/assert"

var key Key = NewKey()
var _ sql.Scanner = &key
var _ driver.Valuer = &key

func TestNewKey(t *testing.T) {
	result := NewKey()
	a.NotNil(t, result, "result is not nil.")
}

func TestScan(t *testing.T) {
	var key Key

	test := func(src interface{}) {
		key = Key("")
		e := key.Scan(src)
		a.NoError(t, e)
	}

	test("string")
	test([]byte("string"))
}

func TestValue(t *testing.T) {
	k := NewKey()
	v, e := k.Value()
	a.NoError(t, e)
	a.Equal(t, v, driver.Value(string(k)), "wrong driver value generated.")
}
