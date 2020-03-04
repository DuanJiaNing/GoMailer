package app

import (
	"net/http"

	netcontext "golang.org/x/net/context"

	"GoMailer/log"
)

var (
	contextKey = "holds a internal *context"
)

type context struct {
	req *http.Request

	outCode   int
	outHeader http.Header
	outBody   []byte
}

func fromContext(ctx netcontext.Context) *context {
	c, _ := ctx.Value(&contextKey).(*context)
	return c
}

func withContext(parent netcontext.Context, c *context) netcontext.Context {
	ctx := netcontext.WithValue(parent, &contextKey, c)
	return ctx
}

func (c *context) Header() http.Header { return c.outHeader }

func (c *context) Write(b []byte) (int, error) {
	if c.outCode == 0 {
		c.WriteHeader(http.StatusOK)
	}
	if len(b) > 0 && !bodyAllowedForStatus(c.outCode) {
		return 0, http.ErrBodyNotAllowed
	}
	c.outBody = append(c.outBody, b...)
	return len(b), nil
}

func (c *context) WriteHeader(code int) {
	if c.outCode != 0 {
		log.Error("WriteHeader called multiple times on request.")
		return
	}
	c.outCode = code
}

// Copied from $GOROOT/src/pkg/net/http/transfer.go. Some response status
// codes do not permit a response body (nor response entity headers such as
// Content-Length, Content-Type, etc).
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == 204:
		return false
	case status == 304:
		return false
	}
	return true
}
