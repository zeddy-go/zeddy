package container

import (
	"reflect"
	"runtime"
	"testing"
)

//func TestContainer_Bind(t *testing.T) {
//	a := 1
//	f := func() {
//		a++
//	}
//	c := &Container{}
//	c.Bind(f)
//}

func a() {
	println(123)
}

func TestXxx(t *testing.T) {

	b := func() {
		println(123)
	}

	c := func() {
		b()
	}

	f := runtime.FuncForPC(reflect.ValueOf(c).Pointer())
	println(f.Name())
}
