package apify

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Params extracts parameters from the request to struct
func Params(r *http.Request, in interface{}) error {
	vars := mux.Vars(r)

	inElem := reflect.ValueOf(in).Elem()
	inType := inElem.Type()

	for i := 0; i < inType.NumField(); i++ {
		field := ParamsStructField(inType.Field(i))

		if value, ok := vars[field.Key()]; ok {
			fieldValue := inElem.Field(i)

			switch kind := fieldValue.Kind(); kind {
			case reflect.String:
				fieldValue.SetString(value)
			case reflect.Int64:
				intValue, err := strconv.ParseInt(value, 10, 64)

				if err != nil {
					return fmt.Errorf("invalid field value format: %w", err)
				}

				fieldValue.SetInt(intValue)
			default:
				return fmt.Errorf("unsupported field kind %s", kind)
			}
		}
	}

	return nil
}

// ParamsStructField wraps reflect.StructField
type ParamsStructField reflect.StructField

// Key returns parameter key for the specific struct field
func (f ParamsStructField) Key() string {
	if val, ok := f.Tag.Lookup("params"); ok {
		return val
	}

	return strings.ToLower(f.Name)
}
