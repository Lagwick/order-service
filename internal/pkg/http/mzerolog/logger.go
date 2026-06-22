package mzerolog

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Lagwick/order-service/internal/pkg/http/httph"
)

type middleware struct {
	log zerolog.Logger

	fromOptions struct {
		skipper func(r *http.Request) bool
	}
}

func (m *middleware) Callback(c *gin.Context) {
	start := time.Now()
	c.Next()
	r := c.Request
	err := httph.ErrorGet(r)
	if m.fromOptions.skipper(r) {
		return
	}
	var builder strings.Builder
	builder.Grow(48 + len(r.RequestURI))

	builder.WriteString(r.Method)
	builder.WriteByte(' ')
	builder.WriteString(r.RequestURI)

	if err != nil {
		builder.WriteString(" finished (or aborted) with error")
	} else {
		builder.WriteString(" finished with no error")
	}

	event := m.log.Debug()
	if err != nil {
		event = m.log.Error().Err(err)
	}

	event = event.
		Dur("exec_time", time.Since(start)).
		Str("client_ip", c.ClientIP())

	if ctxErr := r.Context().Err(); ctxErr != nil {
		event = event.Err(ctxErr)
	}

	event.Msg(builder.String())
}

func NewMiddleware(opts ...Option) gin.HandlerFunc {
	m := &middleware{
		log: log.Logger,
	}

	m.fromOptions.skipper = defaultSkipper

	for _, opt := range opts {
		opt(m)
	}

	return m.Callback
}

func defaultSkipper(_ *http.Request) bool {
	return false
}
