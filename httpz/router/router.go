package router

import "net/http"

type Router interface {
	http.Handler
	Handle(method, path string, handle http.Handler) error
}
