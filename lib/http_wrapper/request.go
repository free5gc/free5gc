//go:binary-only-package

package http_wrapper

import (
	"net/http"
	"net/url"
)

type Request struct {
	Params map[string]string
	Header http.Header
	Query  url.Values
	Body   interface{}
}

func NewRequest(req *http.Request, body interface{}) (ret *Request) {}
