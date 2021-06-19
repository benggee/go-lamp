package httpz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

const (
	// ApplicationJson means application/json.
	ApplicationJson = "application/json"
	// ContentType means Content-Type.
	ContentType = "Content-Type"
)


var (
	errorHandler func(error) (int, interface{})
	lock         sync.RWMutex
)

func Error(w http.ResponseWriter, err error) {
	lock.RLock()
	handler := errorHandler
	lock.RUnlock()

	if handler == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code, body := errorHandler(err)
	e, ok := body.(error)
	if ok {
		http.Error(w, e.Error(), code)
	} else {
		WriteJson(w, code, body)
	}
}

// Success writes HTTP 200 OK into w.
func Success(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// SuccessJson writes v into w with 200 OK.
func SuccessJson(w http.ResponseWriter, v interface{}) {
	WriteJson(w, http.StatusOK, v)
}

// SetErrorHandler sets the error handler, which is called on calling Error.
func SetErrorHandler(handler func(error) (int, interface{})) {
	lock.Lock()
	defer lock.Unlock()
	errorHandler = handler
}

// WriteJson writes v as json string into w with code.
func WriteJson(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(code)

	if bs, err := json.Marshal(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if n, err := w.Write(bs); err != nil {
		// http.ErrHandlerTimeout has been handled by http.TimeoutHandler,
		// so it's ignored here.
		if err != http.ErrHandlerTimeout {
			fmt.Errorf("write response failed, error: %s", err)
		}
	} else if n < len(bs) {
		fmt.Errorf("actual bytes: %d, written bytes: %d", len(bs), n)
	}
}