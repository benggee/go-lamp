package httpz

import (
	"errors"
	"github.com/seepre/go-lamp/httpz/router"
	"log"
	"net/http"
)

type (
	options struct {
		start func(e *engine) error
	}

	RunOption func(*Serve)

	Serve struct {
		engine *engine
		opts options
	}
)

func MustNewServe(c HttpConf, opts ...RunOption) *Serve {
	srv, err := NewServe(c, opts...)
	if err != nil {
		log.Fatal(err)
	}
	return srv
}

func NewServe(c HttpConf, opts ...RunOption) (*Serve, error) {
	if len(opts) > 1 {
		return nil, errors.New("only one run option is allowed")
	}

	serve := &Serve{
		engine: newCore(c),
		opts: options{
			start: func(e *engine) error {
				return e.Start()
			},
		},
	}

	for _, opt := range opts {
		opt(serve)
	}

	return serve, nil
}

func (s *Serve) AddRoutes(rs []router.Route, opts ...RouteOption) *Serve {
	r := router.Routers{
		Routes: rs,
	}

	for _, opt := range opts {
		opt(&r)
	}
	s.engine.AddRoute(r)
	return s
}

func (s *Serve) AddRoute(r router.Route, opts ...RouteOption) *Serve {
	s.AddRoutes([]router.Route{r}, opts...)
	return s
}

func (s *Serve) WithMiddleware(middleware Middleware) {
	s.engine.withMiddleware(middleware)
}

func (s *Serve) WithMiddlewares(middlewares  ...Middleware) {
	for _, m := range middlewares {
		s.WithMiddleware(m)
	}
}

func (s *Serve) Run() {
	handlerError(s.opts.start(s.engine))
}

func (s *Serve) Stop() {
	// TODO something
}

func handlerError(err error) {
	if err == nil || err == http.ErrServerClosed {
		return
	}
	panic(err)
}
