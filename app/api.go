package app

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"

	"GoMailer/log"
)

func Handler() http.HandlerFunc {
	return handleHTTP
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	c := &context{
		req:       r,
		outHeader: w.Header(),
	}
	r = r.WithContext(withContext(r.Context(), c))
	c.req = r

	executeRequestSafely(c, r)
	if c.outCode != 0 {
		w.WriteHeader(c.outCode)
	}
	if c.outBody != nil {
		w.Write(c.outBody)
	}
}

func executeRequestSafely(c *context, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			log.Errorf("%s", renderPanic(x))
			c.outCode = 500
		}
	}()

	http.DefaultServeMux.ServeHTTP(c, r)
}

func renderPanic(x interface{}) string {
	buf := make([]byte, 16<<10) // 16 KB should be plenty
	buf = buf[:runtime.Stack(buf, false)]

	const (
		skipStart  = "app.renderPanic"
		skipFrames = 2
	)
	start := bytes.Index(buf, []byte(skipStart))
	p := start
	for i := 0; i < skipFrames*2 && p+1 < len(buf); i++ {
		p = bytes.IndexByte(buf[p+1:], '\n') + p + 1
		if p < 0 {
			break
		}
	}
	if p >= 0 {
		copy(buf[start:], buf[p+1:])
		buf = buf[:len(buf)-(p+1-start)]
	}

	// Add panic heading.
	head := fmt.Sprintf("panic: %v\n\n", x)
	if len(head) > len(buf) {
		// Extremely unlikely to happen.
		return head
	}
	copy(buf[len(head):], buf)
	copy(buf, head)

	return string(buf)
}
