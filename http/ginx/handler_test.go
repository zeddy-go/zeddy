package ginx

import (
	"reflect"
	"testing"
)

func a(m IMeta) IMeta {
	return nil
}

func TestXxx(t *testing.T) {

	mt := reflect.TypeOf(a)

	println(mt.Out(0).Name())
}
