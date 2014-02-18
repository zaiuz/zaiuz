package zaiuz_test

import "testing"
import "database/sql"
import "database/sql/driver"
import . "github.com/zaiuz/zaiuz"
import a "github.com/stretchr/testify/assert"

var hash Hash = Hash("")
var _ sql.Scanner = &hash
var _ driver.Valuer = &hash
var _ driver.Valuer = hash

func TestNewHash(t *testing.T) {
	result := NewHash("the string")
	a.NotEmpty(t, result, "hash is empty for non-empty string.")
}

func TestHash_MatchOriginal(t *testing.T) {
	original := "the string"
	result := NewHash(original)
	a.False(t, result.MatchOriginal("something else"), "hash compares wrong.")
	a.True(t, result.MatchOriginal(original), "hash compares wrong.")
}

func TestHash_Scan(t *testing.T) {
	var hash Hash

	test := func(src interface{}) {
		previous := hash
		hash = Hash("")
		e := hash.Scan(src)
		a.NoError(t, e)
		a.NotEqual(t, hash, previous, "value unchanged.")
	}

	test("string")
	test([]byte("asdf"))
}

func TestHash_Value(t *testing.T) {
	h := NewHash("asdf")
	v, e := h.Value()
	a.NoError(t, e)
	a.True(t, len(v.(string)) > 0, "result is empty.")
	a.Equal(t, byte(v.(string)[0]), byte('$'), "result looks wrong.")
}
