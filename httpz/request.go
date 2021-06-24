package httpz

import (
	"encoding/json"
	"github.com/seepre/go-lamp/httpz/internal/context"
	"github.com/seepre/go-lamp/httpz/internal/unmarshaler"
	"io/ioutil"
	"net/http"
)


const (
	pathTag = "path"
	formTag = "form"
)

var pathUnmarshaler = unmarshaler.NewUnmarshaler(pathTag)
var formUnmarshaler = unmarshaler.NewUnmarshaler(formTag)

func Parse(r *http.Request, val interface{}) error {
	if err := ParseFromPath(r, val); err != nil {
		return err
	}

	if err := ParseFromForm(r, val); err != nil {
		return err
	}

	return ParseFromBody(r, val)
}

func ParseFromPath(r *http.Request, val interface{}) error {
	pathVars, err := context.GetPathVars(r)
	if err != nil {
		return err
	}

	return pathUnmarshaler.Unmarshal(val, pathVars)
}

func ParseFromForm(r *http.Request, val interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	formMap := make(map[string]string, len(r.Form))
	for k, v := range r.Form {
		formMap[k] = v[0]
	}

	if len(formMap) == 0 {
		return nil
	}

	return formUnmarshaler.Unmarshal(val, formMap)
}

func ParseFromBody(r *http.Request, val interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return nil
	}

	return json.Unmarshal(body, val)
}