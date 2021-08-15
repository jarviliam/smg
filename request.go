package smg

import (
	"net/http"
	"net/url"
)

type Request struct {
	URL     *url.URL
	Headers *http.Header
	Depth   int
	Ctx     *Context
	spider  *Spider
}
