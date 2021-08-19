package spider

import (
	"net/http"

	"github.com/jarviliam/smg/internal/context"
)

type Response struct {
	Code    int
	Body    []byte
	Ctx     *context.Context
	Request *http.Request
}
