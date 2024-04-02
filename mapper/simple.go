package mapper

import (
	"github.com/zeddy-go/zeddy/reflectx"
	"reflect"
)

func SimpleMap(dst any, src any, sets ...func(*Simple)) (err error) {
	m := NewSimple(sets...)
	return m.Map(dst, src)
}

func SimpleMapSlice(dst any, src any, sets ...func(*Simple)) (err error) {
	m := NewSimple(sets...)
	return m.MapSlice(dst, src)
}

func MustSimpleMap(dst any, src any, sets ...func(*Simple)) {
	err := SimpleMap(dst, src, sets...)
	if err != nil {
		panic(err)
	}
}

func MustSimpleMapSlice(dst any, src any, sets ...func(*Simple)) {
	err := SimpleMapSlice(dst, src, sets...)
	if err != nil {
		panic(err)
	}
}

func WithFieldSimple(f func(srcField string) (dstField string)) func(*Simple) {
	return func(mapper *Simple) {
		mapper.fieldMapper = f
	}
}

func NewSimple(sets ...func(*Simple)) *Simple {
	m := &Simple{}
	for _, set := range sets {
		set(m)
	}

	return m
}

type Simple struct {
	fieldMapper func(srcField string) (dstField string)
}

func (s *Simple) MapSlice(dst any, src any) (err error) {
	dstValue := reflect.ValueOf(dst).Elem()
	srcValue := reflect.ValueOf(src)
	for srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	if dstValue.IsNil() {
		dstValue.Set(reflect.MakeSlice(dstValue.Type(), 0, srcValue.Len()))
	}

	itemType := dstValue.Type().Elem()
	for i := 0; i < srcValue.Len(); i++ {
		dstItem := reflect.New(itemType).Elem()
		err = s.mapValue(dstItem, srcValue.Index(i))
		if err != nil {
			return
		}
		dstValue.Set(reflect.Append(dstValue, dstItem))
	}

	return
}

func (s *Simple) Map(dst any, src any) (err error) {
	dstValue := reflect.ValueOf(dst).Elem()

	srcValue := reflect.ValueOf(src)
	for srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	return s.mapValue(dstValue, srcValue)
}

func (s *Simple) mapValue(dstValue reflect.Value, srcValue reflect.Value) (err error) {
	if srcValue.Kind() == reflect.Struct && dstValue.Kind() == reflect.Struct {
		err = s.mapStructToStruct(dstValue, srcValue)
	} else if srcValue.Kind() == reflect.Map && dstValue.Kind() == reflect.Struct {
		err = s.mapMapToStruct(dstValue, srcValue)
	} else if srcValue.Kind() == reflect.Struct && dstValue.Kind() == reflect.Map {
		err = s.mapStructToMap(dstValue, srcValue)
	} else {
		panic("unsupported type")
	}

	return
}

func (s *Simple) mapStructToMap(dst reflect.Value, src reflect.Value) (err error) {
	if dst.IsNil() {
		dst.Set(reflect.MakeMapWithSize(dst.Type(), src.NumField()))
	}
	for i := 0; i < src.NumField(); i++ {
		if src.Type().Field(i).Anonymous {
			err = s.mapStructToMap(dst, src.Field(i))
			if err != nil {
				return
			}
		}

		var fieldName string
		if s.fieldMapper == nil {
			fieldName = src.Type().Field(i).Name
		} else {
			fieldName = s.fieldMapper(src.Type().Field(i).Name)
		}

		dst.SetMapIndex(reflect.ValueOf(fieldName), src.Field(i))
	}

	return
}

func (s *Simple) mapStructToStruct(dst reflect.Value, src reflect.Value) (err error) {
	for i := 0; i < src.NumField(); i++ {
		if src.Type().Field(i).Anonymous {
			err = s.mapStructToStruct(dst, src.Field(i))
			if err != nil {
				return
			}
		}

		var fieldName string
		if s.fieldMapper == nil {
			fieldName = src.Type().Field(i).Name
		} else {
			fieldName = s.fieldMapper(src.Type().Field(i).Name)
		}

		field := dst.FieldByName(fieldName)
		if !field.CanSet() {
			continue
		}

		err = reflectx.SetValue(src.Field(i), field)
		if err != nil {
			return
		}
	}

	return
}

func (s *Simple) mapMapToStruct(dst reflect.Value, src reflect.Value) (err error) {
	iter := src.MapRange()
	for iter.Next() {
		var fieldName string
		if s.fieldMapper == nil {
			fieldName = iter.Key().String()
		} else {
			fieldName = s.fieldMapper(iter.Key().String())
		}

		field := dst.FieldByName(fieldName)
		if !field.CanSet() {
			continue
		}

		err = reflectx.SetValue(iter.Value(), field)
		if err != nil {
			return
		}
	}

	return
}
