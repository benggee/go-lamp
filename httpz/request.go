package httpz

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/seepre/go-lamp/httpz/internal/context"
)

// TODO concurrency problem
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

	rv := reflect.ValueOf(val)
	if err = validatePtr(&rv); err != nil {
		return err
	}

	rtv := reflect.TypeOf(val).Elem()
	if rtv.Kind() != reflect.Struct {
		return errors.New("not valid struct")
	}

	rev := rv.Elem()
	fieldNums := rtv.NumField()
	for i := 0; i < fieldNums; i++ {
		pathTag, ok :=  rtv.Field(i).Tag.Lookup("path")
		if !ok {
			continue
		}
		if _, ok := pathVars[pathTag]; !ok {
			continue
		}

		if rev.Field(i).Kind() != reflect.TypeOf(pathVars[pathTag]).Kind() {
			return  errors.New("data type invalid")
		}

		rev.Field(i).Set(reflect.ValueOf(pathVars[pathTag]))
	}
	return nil
}

func ParseFromForm(r *http.Request, val interface{}) error {
	pathVars, err := context.GetPathVars(r)
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(val)
	if err = validatePtr(&rv); err != nil {
		return err
	}

	rtv := reflect.TypeOf(val).Elem()
	if rtv.Kind() != reflect.Struct {
		return errors.New("not valid struct")
	}

	rev := rv.Elem()
	fieldNums := rtv.NumField()
	for i := 0; i < fieldNums; i++ {
		pathTag, ok :=  rtv.Field(i).Tag.Lookup("form")
		if !ok {
			continue
		}
		if _, ok := pathVars[pathTag]; !ok {
			continue
		}

		if rev.Field(i).Kind() != reflect.TypeOf(pathVars[pathTag]).Kind() {
			return  errors.New("data type invalid")
		}

		rev.Field(i).Set(reflect.ValueOf(pathVars[pathTag]))
	}
	return nil
}

func ParseFromBody(r *http.Request, val interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, val)
}

func validatePtr(p *reflect.Value) error {
	if !p.IsValid() || p.Kind() != reflect.Ptr || p.IsNil() {
		return errors.New("not valid ptr")
	}
	return nil
}
