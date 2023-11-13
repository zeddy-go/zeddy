package container

import "reflect"

func Resolve[T any](key ...string) (result T, err error) {
	var tmp T
	res, err := Default.Resolve(reflect.TypeOf(tmp), key...)
	if err != nil {
		return
	}
	result = res.(T)
	return
}

func Invoke(f any) (results []reflect.Value, err error) {
	return Default.Invoke(f)
}
