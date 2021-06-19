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

func (s *Serve) AddRoutes(rs []router.Route, opts ...RouteOption) {
	r := router.Routers{
		Routes: rs,
	}

	for _, opt := range opts {
		opt(&r)
	}
	s.engine.AddRoute(r)
}

func (s *Serve) AddRoute(r router.Route, opts ...RouteOption) {
	s.AddRoutes([]router.Route{r}, opts...)
}

func (s *Serve) UseMiddleware(middleware Middleware) {
	s.engine.useMiddleware(middleware)
}

func (s *Serve) Run() {
	handlerError(s.opts.start(s.engine))
}

func (s *Serve) Stop() {
	// TODO something
}

func WithNotAllowedHandler(handler http.Handler) RunOption {
	rt := router.NewRouter()
	return WithRouter(rt)
}

func WithRouter(router router.Router) RunOption {
	return func(serve *Serve) {
		serve.opts.start = func(e *engine) error {
			return e.WithRouter(router)
		}
	}
}

func handlerError(err error) {
	if err == nil || err == http.ErrServerClosed {
		return
	}
	panic(err)
}

func WithMiddleware(middleware Middleware, rt ...router.Route) []router.Route {
	routes := make([]router.Route, len(rt))

	for i := range rt {
		route := rt[i]
		routes[i] = router.Route{
			Method:  route.Method,
			Path:    route.Path,
			Handler: middleware(route.Handler),
		}
	}

	return routes
}

func WithMiddlewares(ms []Middleware, rt ...router.Route) []router.Route {
	for i := len(ms) - 1; i >= 0; i-- {
		rt = WithMiddleware(ms[i], rt...)
	}

	return rt
}
