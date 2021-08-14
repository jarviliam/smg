package smg

import (
	"net/http"
	"net/url"
)

type Request struct {
	URL     *url.URL
	Headers *http.Header
	Ctx     *Context
	spider  *Spider
}
