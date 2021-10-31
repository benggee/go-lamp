package mysqlz

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const tag = "db"

type Session struct {
	db *sql.DB
}

func NewSession(c Conf) *Session {
	if c.ConnMaxLifetime == 0{
		c.ConnMaxLifetime = time.Minute * 3
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 3
	}
	if c.MaxOpenConns == 0 {
		c.MaxIdleConns = 10
	}

	opts := []Option{
		setConnMaxLifetime(c.ConnMaxLifetime),
		setMaxIdleConns(c.MaxIdleConns),
		setMaxOpenConns(c.MaxOpenConns),
	}

	return &Session{
		db: newMysql(c.Dns, opts...),
	}
}

func (s *Session) Row(v interface{}, query string, args... interface{}) error {
	rows, err := s.db.Query(query, args)
	if err != nil {
		return err
	}
	return deserializeRow(v, rows)
}

func (s *Session) Rows(v interface{}, query string, args... interface{}) error {
	rows, err := s.db.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return err
	}
	return deserializeRows(v, rows)
}

func (s *Session) Exec(query string, args... interface{}) (sql.Result, error) {
	result, err := s.db.Exec(query, args)
	if err != nil {
		return nil, err
	}
	return result, nil
}


func deserializeRows(v interface{}, rows *sql.Rows) error {
	rv := reflect.ValueOf(v)
	if err := isPtrInvalid(&rv); err != nil {
		return err
	}

	rt := reflect.TypeOf(v)
	rte := rt.Elem()
	rve := rv.Elem()

	if !rve.CanSet() {
		return errors.New("the data can not set")
	}

	// rows must be slice
	if rte.Kind() != reflect.Slice {
		return errors.New("data is not a slice")
	}

	tmpBase := getElem(rte.Elem())

	switch tmpBase.Kind() {
	case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String:
		for rows.Next() {
			value := reflect.New(tmpBase)
			if !rve.CanSet() {
				return errors.New("base type cant not set")
			}
			if err := rows.Scan(value.Interface()); err != nil {
				return err
			}
			if rte.Elem().Kind() == reflect.Ptr {
				rve.Set(reflect.Append(rve, reflect.ValueOf(value)))
			} else {
				rve.Set(reflect.Append(rve, reflect.Indirect(reflect.ValueOf(value))))
			}
		}
	case reflect.Struct:
		columns, err := rows.Columns()
		if err != nil {
			return err
		}

		for rows.Next() {
			value := reflect.New(tmpBase)
			values, err := convertStructMapToSlice(value, columns)
			if err != nil {
				return err
			}

			if err := rows.Scan(values...); err != nil {
				return err
			}

			if rte.Elem().Kind() == reflect.Ptr {
				rve.Set(reflect.Append(rve, value))
			} else {
				rve.Set(reflect.Append(rve, reflect.Indirect(value)))
			}
		}
	default:
		return errors.New("type invalid")
	}

	return nil
}


func deserializeRow(v interface{}, rows *sql.Rows) error {
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return errors.New("ErrNotFound")
	}

	rv := reflect.ValueOf(v)
	if err := isPtrInvalid(&rv); err != nil {
		return err
	}

	rte := reflect.TypeOf(v).Elem()
	rve := rv.Elem()
	switch rte.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		if rve.CanSet() {
			return rows.Scan(v)
		} else {
			return errors.New("not settable")
		}
	case reflect.Struct:
		columns, err := rows.Columns()
		if err != nil {
			return err
		}
		if values, err := convertStructMapToSlice(rve, columns); err != nil {
			return err
		} else {
			return rows.Scan(values...)
		}
	default:
		return errors.New("unsupported value type")
	}
}



func convertStructMapToSlice(v reflect.Value, cols []string) ([]interface{}, error) {
	tagValMap, err := buildTagToValMap(v)
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(cols))
	if len(tagValMap) != 0 {
		for idx, col := range cols {
			if tag, ok := tagValMap[col]; ok {
				values[idx] = tag
			} else {
				var tmpInterface interface{}
				values[idx] = &tmpInterface
			}
		}
		return values, nil
	} else {
		fields := parseFields(v)
		for i := 0; i < len(cols); i++ {
			vF := fields[i]
			if vF.Kind() != reflect.Ptr {
				if !vF.CanAddr() || !vF.Addr().CanInterface() {
					return nil, errors.New("not readable value")
				}
				values[i] = vF.Addr().Interface()
			} else {
				if !vF.CanInterface() {
					return nil, errors.New("nto readable value")
				}
				if vF.IsNil() {
					bVType := getElem(vF.Type())
					vF.Set(reflect.New(bVType))
				}
				values[i] = vF.Interface()
			}
		}
	}

	return values, nil
}


func buildTagToValMap(v reflect.Value) (map[string]interface{}, error) {
	rt := getElem(v.Type())
	size := rt.NumField()
	result := make(map[string]interface{}, size)

	for i := 0; i < size; i++ {
		key := parseTagName(rt.Field(i))
		if len(key) == 0 {
			return nil, nil
		}

		valueField := reflect.Indirect(v).Field(i)
		switch valueField.Kind() {
		case reflect.Ptr:
			if !valueField.CanInterface() {
				return nil, errors.New("not readable value")
			}
			if valueField.IsNil() {
				baseValueType := getElem(valueField.Type())
				valueField.Set(reflect.New(baseValueType))
			}
			result[key] = valueField.Interface()
		default:
			if !valueField.CanAddr() || !valueField.Addr().CanInterface() {
				return nil, errors.New("not readable value")
			}
			result[key] = valueField.Addr().Interface()
		}
	}

	return result, nil
}

func parseFields(v reflect.Value) []reflect.Value {
	var fields []reflect.Value
	indirect := reflect.Indirect(v)

	for i := 0; i < indirect.NumField(); i++ {
		child := indirect.Field(i)
		if child.Kind() == reflect.Ptr && child.IsNil() {
			baseValueType := getElem(child.Type())
			child.Set(reflect.New(baseValueType))
		}

		child = reflect.Indirect(child)
		childType := indirect.Type().Field(i)
		if child.Kind() == reflect.Struct && childType.Anonymous {
			fields = append(fields, parseFields(child)...)
		} else {
			fields = append(fields, child)
		}
	}

	return fields
}

func parseTagName(field reflect.StructField) string {
	key := field.Tag.Get(tag)
	if len(key) == 0 {
		return ""
	} else {
		options := strings.Split(key, ",")
		return options[0]
	}
}


func getElem(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}


func isPtrInvalid(v *reflect.Value) error {
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("pointer is invalid")
	}
	return nil
}