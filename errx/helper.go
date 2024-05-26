package errx

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"google.golang.org/grpc/status"
)

// PrintGrpcStatusErrorDetail 打印err的详细信息
func PrintGrpcStatusErrorDetail(err error) {
	s, ok := status.FromError(err)
	if !ok {
		println("none status error")
		return
	}

	fmt.Printf("%#v\n", NewFromStatus(s))
}

func PbMap2MapStrAny(pbMap *Map) (m map[string]any) {
	m = make(map[string]any, len(pbMap.Fields))
	for key, item := range pbMap.Fields {
		m[key] = PbValue2Any(item)
	}

	return
}

func PbValue2Any(pbValue *Value) (v any) {
	switch x := pbValue.Kind.(type) {
	case *Value_Int32Val:
		v = x.Int32Val
	case *Value_Int64Val:
		v = x.Int64Val
	case *Value_MapVal:
		v = PbMap2MapStrAny(x.MapVal)
	case *Value_ListVal:
		v = PbList2SliceAny(x.ListVal)
	case *Value_StrVal:
		v = x.StrVal
	default:
		panic("unsupported kind")
	}
	return
}

func PbList2SliceAny(pbList *List) (s []any) {
	s = make([]any, len(pbList.List))
	for _, item := range pbList.List {
		res := PbValue2Any(item)
		s = append(s, res)
	}

	return
}

func Map2Pb(val any) (m map[string]*Value) {
	v := reflect.ValueOf(val)
	m = make(map[string]*Value, v.Len())
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key()
		if key.Kind() != reflect.String {
			panic(errors.New("map key can only be string"))
		}
		v := Any2PbValue(iter.Value().Interface())
		if v != nil {
			m[key.String()] = v
		}
	}

	return
}

func Any2PbValue(val any) (value *Value) {
	value = &Value{}
	switch x := val.(type) {
	case int32:
		value.Kind = &Value_Int32Val{
			Int32Val: x,
		}
	case int64:
		value.Kind = &Value_Int64Val{
			Int64Val: x,
		}
	case string:
		value.Kind = &Value_StrVal{
			StrVal: x,
		}
	case int:
		switch strconv.IntSize {
		case 32:
			value.Kind = &Value_Int32Val{
				Int32Val: int32(x),
			}
		case 64:
			value.Kind = &Value_Int64Val{
				Int64Val: int64(x),
			}
		}
	case *Errx:
		value.Kind = &Value_StrVal{
			StrVal: x.ErrStack(),
		}
	case error:
		value.Kind = &Value_StrVal{
			StrVal: x.Error(),
		}
	default:
		v := reflect.ValueOf(val)
		switch v.Kind() {
		case reflect.Map:
			var m map[string]*Value
			m = Map2Pb(v.Interface())
			value.Kind = &Value_MapVal{
				MapVal: &Map{
					Fields: m,
				},
			}
		case reflect.Slice:
			s := &Value_ListVal{
				ListVal: &List{
					List: make([]*Value, 0, v.Len()),
				},
			}
			for i := 0; i < v.Len(); i++ {
				item := v.Index(i)
				pbValue := Any2PbValue(item.Interface())
				s.ListVal.List = append(s.ListVal.List, pbValue)
			}
			value.Kind = s
		default:
			return nil
		}
	}
	return
}
