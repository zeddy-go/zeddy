package reflectx

import (
	"errors"
	"fmt"
	"reflect"
)

func SetValue(srcValue reflect.Value, dest reflect.Value) (err error) {
	if dest.Kind() != reflect.Ptr || dest.IsNil() {
		if dest.CanAddr() {
			dest = dest.Addr()
		} else {
			return errors.New("dest value must be pointer or CanAddr")
		}
	}

	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	if srcValue.Kind() != BaseKind(dest) {
		srcValue, err = ConvertToKind(srcValue, BaseType(dest).Kind())
		if err != nil {
			return errors.New(fmt.Sprintf(
				"kind [%s] can not convert to kind [%s]: %s",
				BaseType(srcValue).Kind().String(),
				BaseType(dest).Kind().String(),
				err.Error(),
			))
		}
	}

	tp := dest.Elem().Type()
	v := srcValue
	for tp.Kind() == reflect.Ptr {
		n := reflect.New(v.Type())
		n.Elem().Set(v)
		v = n
		tp = tp.Elem()
	}

	dest.Elem().Set(v)

	return
}

func BaseType(v reflect.Value) (t reflect.Type) {
	t = v.Type()
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}

func BaseKind(v reflect.Value) reflect.Kind {
	t := v.Type()
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Kind()
}
