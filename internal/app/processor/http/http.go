package rprocessor

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/Lagwick/order-service/internal/app/config/section"
	rhandler "github.com/Lagwick/order-service/internal/app/handler/http"
	"github.com/Lagwick/order-service/internal/app/processor"
	"github.com/Lagwick/order-service/internal/app/util"
	"github.com/Lagwick/order-service/internal/pkg/http/httph"
	"github.com/Lagwick/order-service/internal/pkg/http/mzerolog"
)

type httpProc struct {
	server *http.Server
	addr   string
}

func NewHTTP(
	hHealth rhandler.Health,
	hOrder rhandler.Order,
	cfg section.ProcessorWebServer,
) processor.Processor {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(adaptRequestMiddleware(httph.NewErrorMiddleware()))
	router.Use(mzerolog.NewMiddleware(
		mzerolog.WithSkipper(util.IsFilteredHttpRoute),
	))
	router.Use(gin.Recovery())

	router.NoRoute(handleNotFound)
	vGenericRegHealthCheck(router, hHealth)

	v1 := router.Group("/v1")
	v1RegOrderHandler(v1, hOrder)

	for _, route := range router.Routes() {
		log.Info().
			Str("method", route.Method).
			Str("path", route.Path).
			Msg("HTTP route registered")
	}

	addr := fmt.Sprintf(":%d", cfg.ListenPort)

	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &httpProc{
		server: server,
		addr:   addr,
	}
}

func (p *httpProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	lc := net.ListenConfig{}

	l, err := lc.Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Msg("listen tcp")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.serve(l)
	}()

	processor.WatchForShutdown(ctx, wg, l)

	processor.WatchForShutdown(ctx, wg, processor.NewCloserContextFunc(
		p.server.Shutdown,
		context.Background(),
		5*time.Second,
	))
}

func adaptRequestMiddleware(m httph.Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler := m(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			c.Request = r
		}))

		handler.ServeHTTP(c.Writer, c.Request)
		c.Next()
	}
}

func (p *httpProc) serve(l net.Listener) {
	if err := p.server.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("serve http")
	}
}
