package slicex

import (
	"errors"
	"github.com/zeddy-go/zeddy/reflectx"
	"reflect"
	"strings"
)

func Find[T any](f func(item T) bool, arr ...T) (result T, ok bool) {
	if len(arr) == 0 {
		return
	}

	for _, item := range arr {
		if f(item) {
			ok = true
			result = item
			return
		}
	}

	return
}

func Filter[T any](f func(item T) bool, arr ...T) (list []T) {
	if len(arr) == 0 {
		return
	}

	list = make([]T, 0, len(arr))
	for _, item := range arr {
		if !f(item) {
			list = append(list, item)
		}
	}

	return
}

func Last[T any](arr ...T) (result T) {
	if len(arr) == 0 {
		return
	}

	return arr[len(arr)-1]
}

func Pop[T any](arr ...T) (newArr []T, result T) {
	if len(arr) == 0 {
		return
	}
	return arr[:len(arr)-1], arr[len(arr)-1]
}

func First[T any](arr ...T) (result T) {
	if len(arr) == 0 {
		return
	}

	return arr[0]
}

// Contains find if target in data
func Contains[T comparable](target any, data []T) (found bool) {
	if data == nil {
		return false
	}
	switch x := target.(type) {
	case []T:
		for _, find := range x {
			for _, item := range data {
				if item == find {
					found = true
					return
				}
			}
		}
	case T:
		for _, item := range data {
			found = item == x
			if found {
				return
			}
		}
	default:
		panic(errors.New("type is mismatched"))
	}

	return
}

func MapBy[T comparable, P any](fieldName string, target []P) (result map[T]P) {
	result = make(map[T]P, len(target))
	for _, item := range target {
		v := reflect.ValueOf(item)
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		field := v.FieldByName(fieldName)
		if field.IsValid() {
			result[field.Interface().(T)] = item
		}
	}

	return
}

func GetFieldSlice[T any](target any, key string) []T {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		panic(errors.New("type is not slice or array"))
	}
	finds := make([]reflect.Value, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		finds = append(finds, v.Index(i))
	}

	keys := strings.Split(key, ".")
	for _, key := range keys {
		for idx, item := range finds {
			field := reflectx.FindField(item, key)
			if field.IsValid() {
				finds[idx] = field
			} else {
				finds[idx] = reflect.Value{}
			}
		}
	}

	result := make([]T, 0, len(finds))
	for _, find := range finds {
		result = append(result, find.Interface().(T))
	}

	return result
}
