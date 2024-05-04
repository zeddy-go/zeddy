package container

import (
	"reflect"
	"sync/atomic"
)

var (
	def atomic.Value
)

func init() {
	Set(NewContainer())
}

func Set(container *Container) {
	def.Store(container)
}

func Default() *Container {
	return def.Load().(*Container)
}

func Bind[T any](providerOrInstance any, sets ...func(*bindOpts)) (err error) {
	return Default().Bind(reflect.TypeOf((*T)(nil)).Elem(), reflect.ValueOf(providerOrInstance), sets...)
}

func Resolve[T any](opts ...func(*resolveOpts)) (result T, err error) {
	res, err := Default().Resolve(reflect.TypeOf((*T)(nil)).Elem(), opts...)
	if err != nil {
		return
	}
	result = res.Interface().(T)
	return
}

func MustResolve[T any]() (result T) {
	result, err := Resolve[T]()
	if err != nil {
		panic(err)
	}
	return
}

func Has[T any]() bool {
	return Default().Has(reflect.TypeOf((*T)(nil)).Elem())
}

func Invoke(f any, opts ...func(*invokeOpts)) (err error) {
	results, err := Default().Invoke(reflect.ValueOf(f), opts...)
	if err != nil {
		return
	}
	if len(results) > 0 {
		e := results[len(results)-1].Interface()
		if e != nil {
			err = e.(error)
		}
	}
	return
}
