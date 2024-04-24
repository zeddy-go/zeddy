package reflectx

import (
	"errors"
	"fmt"
	"reflect"
)

func IsSameBaseType(v1, v2 reflect.Value) (err error) {
	t1 := Indirect(v1).Type()
	t2 := Indirect(v2).Type()
	if t1 != t2 {
		return errors.New(fmt.Sprintf("type <%s> and <%s> is not same", t1.String(), t2.String()))
	}

	return
}

func SetValue(dstValue reflect.Value, srcValue reflect.Value) (err error) {
	if err := IsSameBaseType(dstValue, srcValue); Indirect(dstValue).Kind() != reflect.Interface && err != nil {
		return err
	}

	for srcValue.Kind() == reflect.Pointer {
		srcValue = srcValue.Elem()
	}

	for dstValue.Kind() == reflect.Pointer {
		if dstValue.IsNil() && dstValue.CanSet() {
			dstValue.Set(reflect.New(dstValue.Type().Elem()))
		}
		dstValue = dstValue.Elem()
	}

	dstValue.Set(srcValue)

	return
}

func Indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	return v
}
