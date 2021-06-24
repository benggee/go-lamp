package unmarshaler

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type unmarshaler struct {
	key string
}

func NewUnmarshaler(key string) *unmarshaler {
	return &unmarshaler{
		key: key,
	}
}

func (u *unmarshaler) Unmarshal(v interface{}, m map[string]string) error {
	rv := reflect.ValueOf(v)
	if err := validatePtr(&rv); err != nil {
		return err
	}

	rtv := reflect.TypeOf(v).Elem()
	if rtv.Kind() != reflect.Struct {
		return errors.New("not valid struct")
	}

	rev := rv.Elem()
	fieldNums := rtv.NumField()
	for i := 0; i < fieldNums; i++ {
		tag, ok := rtv.Field(i).Tag.Lookup(u.key)
		if !ok {
			continue
		}
		if _, ok = m[tag]; !ok {
			continue
		}

		fieldType := derStructType(rtv.Field(i))
		if err := setValue(fieldType.Kind(), rev.Field(i), m[tag]); err != nil {
			fmt.Println(err,":", rev.Field(i), fieldType.Kind(), tag, m[tag], fieldType.Kind())
			return err
		}
	}

	return nil
}

func derStructType(filed reflect.StructField) reflect.Type {
	if filed.Type.Kind() == reflect.Ptr {
		return filed.Type.Elem()
	}
	return  filed.Type
}

func setValue(kind reflect.Kind, value reflect.Value, str string) error {
	if !value.CanSet() {
		return errors.New("value can not set")
	}

	v, err := convertType(kind, str)
	if err != nil {
		return err
	}

	return setMatchValue(kind, value, v)
}

func convertType(kind reflect.Kind, str string) (interface{}, error) {
	switch kind {
	case reflect.Bool:
		return str == "1" || strings.ToLower(str) == "true", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("the value %q can not parsed as int", str)
		}
		return intVal, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("the value %q can not parsed as uint", str)
		}
		return uintVal, nil
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, fmt.Errorf("the value %q can not parsed as float", str)
		}
		return floatVal, nil
	case reflect.String:
		return str, nil
	default:
		return nil, errors.New("unknown type")
	}
}

func setMatchValue(kind reflect.Kind, value reflect.Value, v interface{}) error  {
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

func validatePtr(v *reflect.Value) error  {
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("invalid ptr")
	}
	return nil
}