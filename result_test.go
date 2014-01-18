package zaiuz

import "testing"
import "./testutil"
import a "github.com/stretchr/testify/assert"

func TestResultFunc(t *testing.T) {
	response, request := testutil.NewTestRequestPair()
	context := NewContext(response, request)

	called := false
	execute := func(c *Context) error {
		a.Equal(t, context, c, "context instance not the given instance.")
		called = true
		return nil
	}

	result := ResultFunc(execute)
	a.NotNil(t, result, "cannot create result from a function.")

	e := result.Execute(context)
	a.NoError(t, e)
	a.True(t, called, "given execute function not called.")
}

func TestDudResult(t *testing.T) {
	result := DudResult()
	a.NotNil(t, result, "cannot create dud result.")

	test := func() {
		e := result.Execute(nil)
		a.NoError(t, e)
	}

	a.NotPanics(t, test, "dud result panics on nil context.")
}

