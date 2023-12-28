package event

import (
	"reflect"
	"testing"
)

type testStruct struct{}

func testF(t testStruct) {
	println("ok")
}
func TestXxx(t *testing.T) {
	t1 := testStruct{}
	reflect.ValueOf(testF).Call([]reflect.Value{reflect.ValueOf(interface{}(t1))})
}
