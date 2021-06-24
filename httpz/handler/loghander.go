package handler

import (
	"github.com/seepre/go-lamp/logz"
	"net/http"
)

func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logz.NewLog(logz.DefaultLogger)
		l.Info(r)
		next.ServeHTTP(w, r)
	})
}
