package router

import (
	"errors"
	"net/http"

	"github.com/seepre/go-lamp/httpz/internal/context"
	"github.com/seepre/go-lamp/httpz/pathz"
)

var (
	ErrInvalidMethod = errors.New("invalid http method")
	ErrInvalidPath   = errors.New("invalid http pathz")
)

type (
	route struct {
		path *pathz.Path
	}
)

func NewRouter() Router {
	return &route{
		path: pathz.NewPath(),
	}
}

func (r *route) Handle(method, routePath string, handler http.Handler) error {
	if !validateMethod(method) {
		return ErrInvalidMethod
	}

	if len(routePath) == 0 || routePath[0] != '/' {
		return ErrInvalidPath
	}

	return r.path.BuildPath(method, routePath, handler)
}

func (r *route) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	handle, paramMap, err := r.path.ParsePath(re.Method, re.URL.Path)
	if err != nil && err != pathz.NotFound {
		http.Error(w, err.Error(), 500)
		return
	}
	if err == pathz.NotFound {
		http.NotFound(w, re)
	}
	if handle == nil {
		http.NotFoundHandler()
		return
	}

	re = context.WithRequestContext(re, paramMap)

	handle.ServeHTTP(w, re)
	return
}

func validateMethod(method string) bool {
	return method == http.MethodDelete || method == http.MethodGet ||
		method == http.MethodHead || method == http.MethodOptions ||
		method == http.MethodPatch || method == http.MethodPost ||
		method == http.MethodPut
}
