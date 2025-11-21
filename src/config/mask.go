package config

import (
	"fmt"
	"reflect"
	"strings"
)

func MaskSensitive[T any](cfg T) T {
	mask(reflect.ValueOf(&cfg).Elem())

	return cfg
}

func mask(in reflect.Value) {
	var struc reflect.Value
	if in.Kind() == reflect.Ptr {
		struc = in.Elem()
	} else {
		struc = in
	}

	for i := 0; i < struc.NumField(); i++ {
		field := struc.Field(i)

		if field.Kind() == reflect.Struct ||
			(field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct) {
			mask(field)
		} else if field.Kind() == reflect.String && struc.Type().Field(i).Tag.Get("sensitive") == "true" {
			raw := *field.Addr().Interface().(*string)

			if len(raw) > 4 {
				field.SetString(fmt.Sprintf("%s%s", strings.Repeat("*", 6), raw[len(raw)-2:]))
			} else {
				field.SetString(strings.Repeat("*", len(raw)))
			}
		}
	}
}
