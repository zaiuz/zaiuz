package testutil

import "io/ioutil"
import "net/url"
import "net/http"
import "net/http/httptest"
import "testing"
import a "github.com/stretchr/testify/assert"

// Method chain context object. Do not use directly, use one of the Http* method to obtain
// this and then use Expect method to setup test expectations.
type ResponseExpectable struct {
	T        *testing.T
	Response *http.Response
	Error    error
}

// Starts a HTTP GET request and returns an object for setting up expectation for the
// result.
func HttpGet(t *testing.T, url string) *ResponseExpectable {
	response, e := http.Get(url)
	return &ResponseExpectable{t, response, e}
}

// Starts a HTTP POST request with the given data payload and returns an object for
// setting up expectation for the result.
func HttpPost(t *testing.T, url string, data url.Values) *ResponseExpectable {
	response, e := http.PostForm(url, data)
	return &ResponseExpectable{t, response, e}
}

// Reads the response body and tests if the response parameters match supplied values.
func (r *ResponseExpectable) Expect(code int, body string) {
	a.NoError(r.T, r.Error, "error while getting response.")
	a.Equal(r.T, r.Response.StatusCode, code, "invalid status code.")

	if len(body) > 0 {
		raw, e := ioutil.ReadAll(r.Response.Body)
		a.NoError(r.T, e, "error while reading response.")
		a.Equal(r.T, string(raw), body, "wrong response body.")
	}
}

// Creates a new test request/response pair for testing against a Context or any
// roundtripping code. The request is a simple GET / request and the response is an
// instance of httptest.ResponseRecorder.
func NewTestRequestPair() (http.ResponseWriter, *http.Request) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	return response, request
}
