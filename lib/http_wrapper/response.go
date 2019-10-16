//go:binary-only-package

package http_wrapper

import (
	"net/http"
)

type Response struct {
	Header http.Header
	Status int
	Body   interface{}
}

func NewResponse(code int, h http.Header, body interface{}) (ret *Response) {}
