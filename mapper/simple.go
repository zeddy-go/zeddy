package mapper

import (
	"errors"
	"github.com/zeddy-go/zeddy/convert"
	"github.com/zeddy-go/zeddy/reflectx"
	"reflect"
	"strings"
)

func SimpleMapSlice(dest any, source any) (err error) {
	dst := reflect.ValueOf(dest)
	src := reflect.ValueOf(source)
	if !isStructSlice(dst, src) {
		return errors.New("is not struct")
	}
	return SimpleMapSliceValue(dst, src)
}

func SimpleMap(dest any, source any) (err error) {
	dst := reflect.ValueOf(dest)
	src := reflectx.BaseValue(reflect.ValueOf(source))
	if reflectx.BaseKind(dst) != reflect.Struct || src.Kind() != reflect.Struct {
		err = errors.New("should be struct")
		return
	}
	return SimpleMapValue(dst, src)
}

func SimpleMapSliceValue(dst, src reflect.Value) (err error) {
	for dst.Kind() == reflect.Pointer {
		dst = dst.Elem()
	}
	if dst.IsNil() {
		dst.Set(reflect.MakeSlice(dst.Type(), 0, 10))
	}
	for i := 0; i < src.Len(); i++ {
		dstItem := reflect.New(dst.Type().Elem()).Elem()
		item := dstItem
		for item.Kind() == reflect.Pointer {
			item.Set(reflect.New(item.Type().Elem()))
			item = item.Elem()
		}
		err = SimpleMapValue(dstItem, reflectx.BaseValue(src.Index(i)))
		if err != nil {
			return
		}
		dst.Set(reflect.Append(dst, dstItem))
	}
	return
}

func SimpleMapValue(dst reflect.Value, src reflect.Value) (err error) {
	for i := 0; i < src.NumField(); i++ {
		var (
			srcField       = reflectx.BaseValue(src.Field(i))
			srcFieldStruct = src.Type().Field(i)
			dstField       reflect.Value
			dstFieldStruct reflect.StructField
		)

		if !srcField.IsValid() || srcField.IsZero() {
			continue
		}
		if srcFieldStruct.Anonymous {
			dstField, dstFieldStruct = findFieldAnonymous(dst, srcField.Type())
			if dstField.IsValid() {
				SimpleMapValue(dstField, srcField)
			} else {
				SimpleMapValue(dst, srcField)
			}
		} else {
			dstField, dstFieldStruct = findFieldByName(dst, srcFieldStruct.Name, false)
			if dstField.IsValid() {
				for dstField.Kind() == reflect.Pointer {
					dstField.Set(reflect.New(dstFieldStruct.Type.Elem()))
					dstField = dstField.Elem()
				}
				if srcField.Kind() == reflect.Struct {
					SimpleMapValue(dstField, srcField)
				} else {
					if srcField.Type() != dstField.Type() {
						var targetSrcField reflect.Value
						targetSrcField, e := convert.ToKindValue(srcField, dstField.Kind())
						if e != nil {
							if isStructSlice(dstField, srcField) {
								SimpleMapSliceValue(dstField, srcField)
							}
							continue
						}
						reflectx.SetValue(dstField, targetSrcField)
					} else {
						reflectx.SetValue(dstField, srcField)
					}
				}
			}
		}
	}

	return
}

func isStructSlice(dst reflect.Value, src reflect.Value) bool {
	if reflectx.BaseKindByType(src.Type()) != reflect.Slice || reflectx.BaseKindByType(dst.Type()) != reflect.Slice {
		return false
	}
	if reflectx.BaseKindByType(reflectx.BaseType(src).Elem()) != reflect.Struct || reflectx.BaseKindByType(reflectx.BaseType(dst).Elem()) != reflect.Struct {
		return false
	}

	return true
}

func findFieldByName(v reflect.Value, name string, caseSensitive bool) (field reflect.Value, fieldStruct reflect.StructField) {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		fs := v.Type().Field(i)
		if fs.Anonymous {
			dstField := v.Field(i)
			for dstField.Kind() == reflect.Pointer {
				dstField.Set(reflect.New(fs.Type.Elem()))
				dstField = dstField.Elem()
			}
			field, fieldStruct = findFieldByName(v.Field(i), name, caseSensitive)
			if field.IsValid() {
				return
			}
		}
		fieldName := fs.Name
		if (caseSensitive && fieldName == name) || (!caseSensitive && strings.ToLower(fieldName) == strings.ToLower(name)) {
			return v.Field(i), fs
		}
	}

	return
}

func findFieldAnonymous(v reflect.Value, fieldType reflect.Type) (field reflect.Value, fieldStruct reflect.StructField) {
	v = reflectx.BaseValue(v)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Type() == fieldType {
			return f, v.Type().Field(i)
		}
	}

	return
}
