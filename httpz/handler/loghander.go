package handler

import (
	"net/http"
)

func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO write log
		next.ServeHTTP(w, r)
	})
}
