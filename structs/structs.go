package structs

import (
	"fmt"
	"reflect"
)

func GetStructTags(v interface{}, tag string) []string {
	val := reflect.ValueOf(v)
	names := []string{}

	if val.Kind() != reflect.Struct {
		if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
			val = val.Elem()
		} else {
			panic(fmt.Sprintf("GetStructTags: only structs are allowed got %s", reflect.TypeOf(v).Name()))
		}
	}
	for i := 0; i < val.Type().NumField(); i++ {
		name := val.Type().Field(i).Tag.Get(tag)
		if name == "" {
			name = val.Type().Field(i).Name
		}

		if name != "-" {
			names = append(names, name)
		}
	}

	return names
}
