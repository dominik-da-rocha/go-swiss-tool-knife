package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log/slog"
	"reflect"
	"sort"
)

func HashSha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	hashSum := hash.Sum(nil)
	return hashSum
}

func NewRandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NewRandomString(length int) string {
	bytes, _ := NewRandomBytes(length)
	return base64.URLEncoding.EncodeToString(bytes)
}

func IsNilOrEmpty(text *string) bool {
	return text == nil || *text == ""
}

func IsAnyOfOrString(text *string, allowed []string, def string) string {
	if IsNilOrEmpty(text) {
		return def
	}
	idx := IndexOfString(allowed, *text)
	if idx < 0 {
		return def
	}
	return *text
}

func OrInt64(value *int64, def int64) int64 {
	if value == nil {
		return def
	} else {
		return *value
	}
}

func OrString(value *string, def string) string {
	if value == nil {
		return def
	} else {
		return *value
	}
}

func IndexOfString(slice []string, target string) int {
	for idx, s := range slice {
		if s == target {
			return idx
		}
	}
	return -1
}

func StructTagNames(v interface{}, tag string) []string {
	val := reflect.TypeOf(v)

	names := []string{}
	for i := 0; i < val.NumField(); i++ {
		name := val.Field(i).Tag.Get(tag)
		if name == "" {
			name = val.Field(i).Name
		}
		names = append(names, name)
	}

	return names
}

func GetStructTags(v interface{}, tag string) []string {
	val := reflect.ValueOf(v)
	names := []string{}

	if val.Kind() != reflect.Struct {
		if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
			val = val.Elem()
		} else {
			slog.Error("GetStructTags: only structs are allowed got", "type", reflect.TypeOf(v).Name())
			panic("GetStructTags")
		}
	}
	for i := 0; i < val.Type().NumField(); i++ {
		name := val.Type().Field(i).Tag.Get(tag)
		if name == "" {
			name = val.Type().Field(i).Name
		} else if name != "-" {
			names = append(names, name)
		}
	}

	return names
}

func Uups(err error) {
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}

func MustBeTrue(cond bool, msg string, args ...interface{}) {
	if !cond {
		slog.Error("msg", args...)
		panic(errors.New(msg))
	}
}

func RemoveFromStings(texts []string, toRemove string) []string {
	stripped := []string{}
	for _, text := range texts {
		if text != toRemove {
			stripped = append(stripped, text)
		}
	}
	return texts
}

func StructFindFirstTag(v interface{}, tagToFind string) string {
	val := reflect.TypeOf(v)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		name, found := field.Tag.Lookup(tagToFind)
		if found {
			return name
		}
	}
	return ""
}

func StructSetValueByTag(v interface{}, tag string, key string, newValue interface{}) bool {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
			val = val.Elem()
		} else {
			slog.Error("StructSetValueByTag: only structs are allowed got", "type", reflect.TypeOf(v).Name())
			panic("StructSetValueByTag")
		}
	}

	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		name, found := field.Tag.Lookup(tag)
		if found && name == key {
			if val.Field(i).CanSet() {
				val.Field(i).Set(reflect.ValueOf(newValue))
			}
			return true
		}
	}
	return false
}

func SortedContainsStrings(list []string, search string) bool {
	idx := sort.SearchStrings(list, search)
	if idx >= len(list) {
		return false
	} else {
		return list[idx] == search
	}
}
