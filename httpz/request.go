package httpz

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/seepre/go-lamp/httpz/internal/context"
)

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
		return errors.New("invalid struct")
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

		var fieldType reflect.Type
		if rtv.Field(i).Type.Kind() == reflect.Ptr {
			fieldType = rtv.Field(i).Type.Elem()
		} else {
			fieldType = rtv.Field(i).Type
		}

		if err = setValue(fieldType.Kind(), rev.Field(i), pathVars[pathTag]); err != nil {
			return err
		}
	}
	return nil
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

	rv := reflect.ValueOf(val)
	if err := validatePtr(&rv); err != nil {
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
		if _, ok = formMap[pathTag]; !ok {
			continue
		}

		var fieldType reflect.Type
		if rtv.Field(i).Type.Kind() == reflect.Ptr {
			fieldType = rtv.Field(i).Type.Elem()
		} else {
			fieldType = rtv.Field(i).Type
		}

		if err := setValue(fieldType.Kind(), rev.Field(i), formMap[pathTag]); err != nil {
			return err
		}
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

func convertType(kind reflect.Kind, str string) (interface{}, error) {
	switch kind {
	case reflect.Bool:
		return str == "1" || strings.ToLower(str) == "true", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("the value %q cannot parsed as int", str)
		}

		return intValue, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("the value %q cannot parsed as uint", str)
		}

		return uintValue, nil
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, fmt.Errorf("the value %q cannot parsed as float", str)
		}

		return floatValue, nil
	case reflect.String:
		return str, nil
	default:
		return nil, errors.New("unknown type")
	}
}


func setMatchedValue(kind reflect.Kind, value reflect.Value, v interface{}) error {
	switch kind {
	case reflect.Bool:
		value.SetBool(v.(bool))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(v.(int64))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value.SetUint(v.(uint64))
	case reflect.Float32, reflect.Float64:
		value.SetFloat(v.(float64))
	case reflect.String:
		value.SetString(v.(string))
	default:
		return errors.New("unknown type")
	}

	return nil
}

func setValue(kind reflect.Kind, value reflect.Value, str string) error {
	if !value.CanSet() {
		return errors.New("value can not set")
	}

	v, err := convertType(kind, str)
	if err != nil {
		return err
	}

	return setMatchedValue(kind, value, v)
}