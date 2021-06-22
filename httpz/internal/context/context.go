package context

import (
	"context"
	"errors"
	"net/http"
)

type contextKey string

func (c contextKey) String() string {
	return "/param/internal/" + string(c)
}

var pathVars = contextKey("pathVars")

func WithRequestContext(r *http.Request, val interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), pathVars, val))
}

func GetPathVars(r *http.Request) (map[string]string, error) {
	varMap, ok := r.Context().Value(pathVars).(map[string]string)
	if !ok {
		return nil, errors.New("path vars invalid")
	}
	return varMap, nil
}
