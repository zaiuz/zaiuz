package zaiuz

import "testing"

// import "./testutil"
import a "github.com/stretchr/testify/assert"

var _ Result = &HttpResult{}

func TestNewHttpResult(t *testing.T) {
	result := NewHttpResult(123, "Content-Type", "text/plain").(*HttpResult)
	a.NotNil(t, result, "cannot create base http result.")
	a.Equal(t, 123, result.Code, "status code not saved.")
	a.NotNil(t, result.Headers, "headers map initialized.")

	contentType, ok := result.Headers["Content-Type"]
	a.True(t, ok, "given header not saved.")
	a.Equal(t, []string{"text/plain"}, contentType, "wrong header value saved.")

	// WAIT: https://github.com/stretchr/testify/issues/34
	// result = NewHttpResult(123, "").(*HttpResult)
	// a.Equal(t, 0, len(result.Headers), "headers list not empty when empty string given.")
}
