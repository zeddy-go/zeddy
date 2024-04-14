package mapper

import (
	"errors"
	"github.com/zeddy-go/zeddy/convert"
	"github.com/zeddy-go/zeddy/reflectx"
	"reflect"
	"strings"
)

func SimpleMap(dest any, source any) (err error) {
	dst := reflect.ValueOf(dest)
	src := reflect.ValueOf(source)
	for src.Kind() == reflect.Ptr {
		src = src.Elem()
	}
	if reflectx.BaseKind(dst) != reflect.Struct || src.Kind() != reflect.Struct {
		err = errors.New("should be struct")
		return
	}
	return SimpleMapValue(dst, src)
}

func SimpleMapValue(dst reflect.Value, src reflect.Value) (err error) {
	for i := 0; i < src.NumField(); i++ {
		var (
			srcField       = src.Field(i)
			srcFieldStruct = src.Type().Field(i)
			dstField       reflect.Value
		)
		dstField = findField(dst, srcFieldStruct, srcField, false)
		if !dstField.IsValid() {
			if srcFieldStruct.Anonymous {
				err = SimpleMapValue(dst, srcField)
				continue
			} else {
				err = errors.New("field is not valid")
				return
			}
		}

		if reflectx.BaseKind(srcField) == reflect.Struct {
			err = SimpleMapValue(dstField, srcField)
			if err != nil {
				return
			}
		} else if reflectx.BaseKind(dstField) != reflectx.BaseKind(srcField) {
			srcField, err = convert.ToKind(srcField, dstField.Kind())
			if err != nil {
				return
			}
		}

		err = reflectx.SetValue(dstField, srcField)
		if err != nil {
			return
		}
	}

	return
}

func findField(v reflect.Value, sf reflect.StructField, f reflect.Value, caseSensitive bool) (r reflect.Value) {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if sf.Anonymous {
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).Anonymous {
				vf := v.Field(i)
				if vf.Type() == f.Type() {
					r = vf
				}
			}
		}
	} else {
		if caseSensitive {
			r = v.FieldByName(sf.Name)
		} else {
			vType := v.Type()
			for i := 0; i < v.NumField(); i++ {
				vFieldStruct := vType.Field(i)
				if vFieldStruct.Anonymous {
					r = findField(v.Field(i), sf, f, caseSensitive)
					if !r.IsValid() {
						continue
					} else {
						break
					}
				} else if strings.ToLower(vFieldStruct.Name) == strings.ToLower(sf.Name) {
					r = v.Field(i)
					break
				}
			}
		}
	}

	return
}
