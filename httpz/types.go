package httpz

import (
	"github.com/seepre/go-lamp/httpz/router"
	"net/http"
)

type (
	Middleware func(next http.HandlerFunc) http.HandlerFunc

	RouteOption func(r *router.Routers)
)
