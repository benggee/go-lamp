package httpz

import (
	"errors"
	"net/http"

	"github.com/seepre/go-lamp/httpz/handler"
	"github.com/seepre/go-lamp/httpz/middleware"
	"github.com/seepre/go-lamp/httpz/router"
)

var ErrSignatureConfig = errors.New("bad config for Signature")

type engine struct {
	conf        HttpConf
	routes      []router.Routers
	middlewares []Middleware
}

func newCore(c HttpConf) *engine {
	co := &engine{
		conf: c,
	}

	return co
}

func (e *engine) AddRoute(r router.Routers) {
	e.routes = append(e.routes, r)
}

func (e *engine) Start() error {
	return e.WithRouter(router.NewRouter())
}

func (e *engine) WithRouter(router router.Router) error {
	if err := e.bindRoutes(router); err != nil {
		return err
	}
	return Start(e.conf.Addr, router)
}

func (e *engine) bindRoutes(router router.Router) error {
	for _, rts := range e.routes {
		for _, route := range rts.Routes {
			if err := e.bindRoute(router, route); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *engine) bindRoute(router router.Router, route router.Route) error {
	w := middleware.New(handler.LogHandler,
		handler.ParseFormHandler,
		handler.RecoverHandler)
	for _, m := range e.middlewares {
		w = w.Append(middleware.ConvertMiddleware(m))
	}
	handle := w.ThenFunc(route.Handler)

	return router.Handle(route.Method, route.Path, handle)
}

func (e *engine) withMiddleware(middleware Middleware) {
	e.middlewares = append(e.middlewares, middleware)
}

func Start(addr string, handle http.Handler) error {
	return start(addr, handle, func(srv *http.Server) error {
		return srv.ListenAndServe()
	})
}

func start(addr string, handle http.Handler, run func(srv *http.Server) error) error {
	server := &http.Server{
		Addr:    addr,
		Handler: handle,
	}

	return run(server)
}
