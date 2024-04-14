package convert

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/zeddy-go/zeddy/reflectx"
	"reflect"
	"testing"
)

func TestXxx(t *testing.T) {
	var i any
	value := reflect.ValueOf(&i)
	println(value.Elem().Type().Kind().String())
	if value.Type().Kind() == reflect.Interface {
		fmt.Println("The reflected object is an interface.")
	} else {
		fmt.Println("The reflected object is not an interface.")
	}
}

func TestSetValue(t *testing.T) {
	t.Run("kind", func(t *testing.T) {
		{
			a := 1
			var b int
			err := reflectx.SetValue(reflect.ValueOf(&b), reflect.ValueOf(a))
			require.NoError(t, err)
			require.Equal(t, a, b)
		}

		{
			a := 1
			var b any
			err := reflectx.SetValue(reflect.ValueOf(&b), reflect.ValueOf(a))
			require.NoError(t, err)
			require.Equal(t, a, b)
		}

		{
			a := 1
			var b *int
			err := reflectx.SetValue(reflect.ValueOf(&b), reflect.ValueOf(a))
			require.NoError(t, err)
			require.Equal(t, a, *b)
		}

		{
			a := 1
			var b **int
			err := reflectx.SetValue(reflect.ValueOf(&b), reflect.ValueOf(a))
			require.NoError(t, err)
			require.Equal(t, a, **b)
		}

		{
			type test struct{ a int }
			a := test{a: 1}
			var b test
			err := reflectx.SetValue(reflect.ValueOf(&b), reflect.ValueOf(a))
			require.NoError(t, err)
			require.Equal(t, a, b)
			require.NotSame(t, &a, &b)
		}

		{
			type test struct{ a int }
			a := test{a: 1}
			var b *test
			err := reflectx.SetValue(reflect.ValueOf(&b), reflect.ValueOf(a))
			require.NoError(t, err)
			require.Equal(t, a, *b)
			require.NotSame(t, &a, b)
		}
	})
}

func TestToKind(t *testing.T) {
	{
		var a any = false
		{
			r, err := ToKind(&a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, false, r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, 0, r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "false", r.Interface())
		}
	}
	{
		a := false
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, false, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, 0, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "false", r.Interface())
		}
	}
	{
		a := true
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, 1, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "true", r.Interface())
		}
	}
	{
		var a any = 1
		{
			r, err := ToKind(&a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := 1
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := 0
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, false, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "0", r.Interface())
		}
	}
	{
		a := int8(1)
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := int16(1)
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := int32(1)
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := int64(1)
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := uint(1)
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		var a any = uint(1)
		{
			r, err := ToKind(&a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := uint8(1)
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := uint16(1)
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := uint32(1)
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := uint64(1)
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		var a any = float32(1)
		{
			r, err := ToKind(&a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(1), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := float32(1)
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := float64(1)
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(a), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := "1"
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, true, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, 1, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(1), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "1", r.Interface())
		}
	}
	{
		a := "0"
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, false, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "0", r.Interface())
		}
	}
	{
		a := "false"
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, false, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "false", r.Interface())
		}
	}
	{
		a := "False"
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, false, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "False", r.Interface())
		}
	}
	{
		a := ""
		{
			r, err := ToKind(a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, false, r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(0), r.Interface())
		}
		{
			r, err := ToKind(a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "", r.Interface())
		}
	}
	{
		var a any = ""
		{
			r, err := ToKind(&a, reflect.Bool)
			require.NoError(t, err)
			require.Equal(t, false, r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int)
			require.NoError(t, err)
			require.Equal(t, int(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int8)
			require.NoError(t, err)
			require.Equal(t, int8(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int16)
			require.NoError(t, err)
			require.Equal(t, int16(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int32)
			require.NoError(t, err)
			require.Equal(t, int32(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Int64)
			require.NoError(t, err)
			require.Equal(t, int64(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint)
			require.NoError(t, err)
			require.Equal(t, uint(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint8)
			require.NoError(t, err)
			require.Equal(t, uint8(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint16)
			require.NoError(t, err)
			require.Equal(t, uint16(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint32)
			require.NoError(t, err)
			require.Equal(t, uint32(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Uint64)
			require.NoError(t, err)
			require.Equal(t, uint64(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Float32)
			require.NoError(t, err)
			require.Equal(t, float32(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.Float64)
			require.NoError(t, err)
			require.Equal(t, float64(0), r.Interface())
		}
		{
			r, err := ToKind(&a, reflect.String)
			require.NoError(t, err)
			require.Equal(t, "", r.Interface())
		}
	}
}
