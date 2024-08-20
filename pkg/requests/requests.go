package requests

import "net/http"

type Request struct {
	Method string
	Url    string
	Body   []byte
	Header http.Header
}
