package testutil

import "io/ioutil"
import "net/url"
import "net/http"
import "net/http/httptest"
import "testing"
import a "github.com/stretchr/testify/assert"

type ResponseExpectable struct {
	T        *testing.T
	Response *http.Response
	Error    error
}

func HttpGet(t *testing.T, url string) *ResponseExpectable {
	response, e := http.Get(url)
	return &ResponseExpectable{t, response, e}
}

func HttpPost(t *testing.T, url string, data url.Values) *ResponseExpectable {
	response, e := http.PostForm(url, data)
	return &ResponseExpectable{t, response, e}
}

func (r *ResponseExpectable) Expect(code int, body string) {
	a.NoError(r.T, r.Error, "error while getting response.")
	a.Equal(r.T, r.Response.StatusCode, code, "invalid status code.")

	if len(body) > 0 {
		raw, e := ioutil.ReadAll(r.Response.Body)
		a.NoError(r.T, e, "error while reading response.")
		a.Equal(r.T, string(raw), body, "wrong response body.")
	}
}

func NewTestRequestPair() (http.ResponseWriter, *http.Request) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	return response, request
}
