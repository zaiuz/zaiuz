package zaiuz_test

import "testing"
import "database/sql"
import "database/sql/driver"
import . "github.com/zaiuz/zaiuz"
import a "github.com/stretchr/testify/assert"

var key Key = NewKey()
var _ sql.Scanner = &key
var _ driver.Valuer = &key
var _ driver.Valuer = key // extra ver for convenience

func TestNewKey(t *testing.T) {
	result := NewKey()
	a.NotNil(t, result, "result is not nil.")
}

func TestKey_Scan(t *testing.T) {
	var key Key

	test := func(src interface{}) {
		key = Key("")
		e := key.Scan(src)
		a.NoError(t, e)
	}

	test("string")
	test([]byte("string"))
}

func TestKey_Value(t *testing.T) {
	k := NewKey()
	v, e := k.Value()
	a.NoError(t, e)
	a.Equal(t, v, driver.Value(string(k)), "wrong driver value generated.")
}
