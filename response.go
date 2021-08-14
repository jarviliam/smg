package smg

import "net/http"

//File for Response Related handling

type Response struct {
	Code    int
	Body    []byte
	Ctx     *Context
	Request *http.Request
}
