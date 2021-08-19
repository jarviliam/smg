package spider

import (
	"net/http"
	"net/url"

	"github.com/jarviliam/smg/internal/context"
)

type Request struct {
	URL     *url.URL
	Headers *http.Header
	Depth   int
	Ctx     *context.Context
	spider  *Spider
}
