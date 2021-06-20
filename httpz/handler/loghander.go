package handler

import (
	"github.com/seepre/go-lamp/log"
	"net/http"
)

func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := log.NewLog(log.DefaultLogger)
		l.Info(r)
		next.ServeHTTP(w, r)
	})
}
