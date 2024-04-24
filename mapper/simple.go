package mapper

import (
	"fmt"
	"github.com/zeddy-go/zeddy/errx"
	"github.com/zeddy-go/zeddy/reflectx"
	"reflect"
	"strings"
)

func SimpleMap(dst any, src any) (err error) {
	return SimpleMapValue(reflect.ValueOf(dst), reflect.ValueOf(src))
}

func SimpleMapValue(dst reflect.Value, src reflect.Value) (err error) {
	src = reflectx.Indirect(src)
	switch src.Kind() {
	case reflect.Struct:
		return SimpleMapStructValueTo(dst, src)
	case reflect.Slice:
		return SimpleMapSliceValueTo(dst, src)
	default:
		if reflectx.Indirect(dst).Type() == src.Type() {
			return reflectx.SetValue(dst, src)
		}
	}

	return
}

func SimpleMapSliceValueTo(dst, src reflect.Value) (err error) {
	dst = reflectx.IndirectAndFill(dst)
	switch dst.Kind() {
	case reflect.Slice:
		return SimpleMapSliceValueToSlice(dst, src)
	default:
		return errx.New(fmt.Sprintf("unsupported src type <%s> to dst type <%s>", src.Type().String(), dst.Type().String()))
	}
}

func SimpleMapSliceValueToSlice(dst reflect.Value, src reflect.Value) (err error) {
	src = reflectx.Indirect(src)
	dst = reflectx.Indirect(dst)
	dstItemType := dst.Type().Elem()
	dst.Set(reflect.MakeSlice(dst.Type(), 0, src.Len()))
	for i := 0; i < src.Len(); i++ {
		dstItem := reflect.New(dstItemType)
		err = SimpleMapValue(dstItem.Elem(), src.Index(i))
		if err != nil {
			return
		}
		dst.Set(reflect.Append(dst, dstItem.Elem()))
	}
	return
}

func SimpleMapStructValueTo(dst reflect.Value, src reflect.Value) (err error) {
	dst = reflectx.IndirectAndFill(dst)
	switch dst.Kind() {
	case reflect.Struct:
		return SimpleMapStructValueToStruct(dst, src)
	default:
		return errx.New(fmt.Sprintf("unsupported src type <%s> to dst type <%s>", src.Type().String(), dst.Type().String()))
	}
}

func SimpleMapStructValueToStruct(dst reflect.Value, src reflect.Value) (err error) {
	dst = reflectx.IndirectAndFill(dst)
	src = reflectx.Indirect(src)
	for i := 0; i < src.NumField(); i++ {
		srcField := src.Field(i)
		srcFieldStruct := src.Type().Field(i)
		if srcFieldStruct.Anonymous {
			dstField := findAnonymous(dst, srcField.Type())
			if !dstField.IsValid() {
				err = SimpleMapStructValueToStruct(dst, srcField)
			} else {
				err = SimpleMapStructValueToStruct(dstField, srcField)
			}
			if err != nil {
				return
			}
			continue
		}

		dstField := findName(dst, srcFieldStruct.Name, false)
		if dstField.IsValid() {
			err = SimpleMapValue(dstField, srcField)
			if err != nil {
				return
			}
		}
	}
	return
}

func findName(v reflect.Value, name string, caseSensitive bool) (field reflect.Value) {
	v = reflectx.IndirectAndFill(v)
	for i := 0; i < v.NumField(); i++ {
		targetStruct := v.Type().Field(i)
		target := v.Field(i)
		if targetStruct.Anonymous {
			f := findName(target, name, caseSensitive)
			if f.IsValid() {
				field = f
				return
			}
		} else {
			switch caseSensitive {
			case true:
				if targetStruct.Name == name {
					field = target
					return
				}
			case false:
				if strings.ToUpper(targetStruct.Name) == strings.ToUpper(name) {
					field = target
					return
				}
			}
		}
	}
	return
}

func findAnonymous(v reflect.Value, t reflect.Type) (field reflect.Value) {
	t = indirectType(t)
	v = reflectx.IndirectAndFill(v)
	for i := 0; i < v.NumField(); i++ {
		if !v.Type().Field(i).Anonymous {
			continue
		}

		f := v.Field(i)
		if indirectType(f.Type()) != t {
			continue
		}

		return f
	}

	return
}

func indirectType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return t
}
