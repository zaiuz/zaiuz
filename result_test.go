package zaiuz

import "net/http"
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

	e := result.Render(context)
	a.NoError(t, e)
	a.True(t, called, "given execute function not called.")
}

func TestDudResult(t *testing.T) {
	result := DudResult()
	a.NotNil(t, result, "cannot create dud result.")

	test := func() {
		e := result.Render(nil)
		a.NoError(t, e)
	}

	a.NotPanics(t, test, "dud result panics on nil context.")
}

func checkResult(t *testing.T, result Result, code int, headers http.Header, body string) {
	response, request := testutil.NewTestRequestPair()
	context := NewContext(response, request)

	e := result.Render(context)
	a.NoError(t, e)

	recorder := response.(*httptest.ResponseRecorder)
	headers := recorder.Header()

	a.Equal(t, code, recorder.Code, "status code incorrect.")
	for k, v := range headers {
		value, ok := headers[k]
		a.True(t, ok, "result does not contains key `%s`.", k)
		a.Equal(t, value, v, "header `%s` have incorrect value.", k)
	}

	if body != "" {
		a.Equal(t, body, string(recorder.Body.Bytes()), "body mismatch.")
	}
}
