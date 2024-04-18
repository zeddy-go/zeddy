package reflectx

import (
	"errors"
	"fmt"
	"reflect"
)

func GetTag(field reflect.StructField, tags ...string) (content string) {
	for content == "" && len(tags) > 0 {
		tag := tags[0]
		tags = tags[1:]
		content = field.Tag.Get(tag)
	}

	return
}

func BaseValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func IsSameBaseType(v1, v2 reflect.Value) (err error) {
	t1 := BaseType(v1)
	t2 := BaseType(v2)
	if t1 != t2 {
		return errors.New(fmt.Sprintf("type <%s> and <%s> is not same", t1.String(), t2.String()))
	}

	return
}

func SetValue(dstValue reflect.Value, srcValue reflect.Value) (err error) {
	if err := IsSameBaseType(dstValue, srcValue); BaseKind(dstValue) != reflect.Interface && err != nil {
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
