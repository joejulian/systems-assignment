package flatfile

import (
	"reflect"
)

// assignBasedOnKind will assign the field data to the field based on the field kind for supported kinds
func assignBasedOnKind(kind reflect.Kind, field reflect.Value, fieldData []byte, tag *flatfileTag) error {
	switch kind {
	case reflect.Pointer:
		// Prevent infinite loop if field is pointing to itself
		if field.Elem().Kind() == reflect.Interface && field.Elem().Elem() == field {
			return nil
		}
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		if field.Type().NumMethod() > 0 && field.CanInterface() {
			if u, ok := field.Interface().(Unmarshaler); ok {
				return u.UnmarshalText(fieldData)
			}
		}
	case reflect.String:
		field.Set(reflect.ValueOf(string(fieldData)))
	}
	return nil
}
