package middleware

import "net/http"

type Constructor func(handler http.Handler) http.Handler

type Ware struct {
	constructors []Constructor
}

func New(constructors ...Constructor) Ware {
	return Ware{append(([]Constructor)(nil), constructors...)}
}

func (w Ware) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := range w.constructors {
		h = w.constructors[len(w.constructors)-1-i](h)
	}

	return h
}

func (w Ware) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		w.Then(nil)
	}
	return w.Then(fn)
}

func (w Ware) Append(constructors ...Constructor) Ware {
	newCons := make([]Constructor, 0, len(w.constructors)+len(constructors))
	newCons = append(newCons, w.constructors...)
	newCons = append(newCons, constructors...)

	return Ware{newCons}
}

func (w Ware) Extend(ware Ware) Ware {
	return w.Append(ware.constructors...)
}

func ConvertMiddleware(middleware func(http.HandlerFunc) http.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware(next.ServeHTTP)
	}
}
