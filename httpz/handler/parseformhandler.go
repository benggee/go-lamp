package handler

import (
	"net/http"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

func ParseFormHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.ParseMultipartForm(defaultMaxMemory)

		next.ServeHTTP(w, r)
	})
}
