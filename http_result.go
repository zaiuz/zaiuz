package zaiuz

import "net/http"

type HttpResult struct {
	Code    int
	Headers http.Header
}

func NewHttpResult(code int, header string, values ...string) Result {
	headers := make(http.Header)
	if header != "" {
		headers[header] = values
	}

	return &HttpResult{code, headers}
}

func (r *HttpResult) Execute(c *Context) error {
	return nil
}
