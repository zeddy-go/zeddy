package container

import (
	"reflect"
)

var Default = NewContainer()

func Bind[T any](providerOrInstance any, sets ...func(*Stuff)) {
	Default.Bind(reflect.ValueOf(new(T)).Elem().Type(), NewStuff(providerOrInstance, sets...))
}

func Register(providerOrInstance any, sets ...func(*Stuff)) {
	Default.Register(NewStuff(providerOrInstance, sets...))
}

func Resolve[T any](key ...string) (result T, err error) {
	tmp := new(T)
	res, err := Default.Resolve(reflect.ValueOf(tmp).Elem().Type(), key...)
	if err != nil {
		return
	}
	result = res.(T)
	return
}

func ResolveType(t reflect.Type, key ...string) (reflect.Value, error) {
	return Default.resolveValue(t, key...)
}

func Invoke(f any) (results []reflect.Value, err error) {
	return Default.Invoke(f)
}
