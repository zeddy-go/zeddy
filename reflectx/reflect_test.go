package reflectx

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestFindFieldIndex(t *testing.T) {
	type Param struct {
		A int
	}

	p := Param{
		A: 1,
	}

	p2 := &Param{
		A: 1,
	}

	{
		result := FindField(reflect.ValueOf(p), "A")
		require.Equal(t, 1, result.Interface())
	}
	{
		result := FindField(reflect.ValueOf(p), "B")
		require.False(t, result.IsValid())
	}
	{
		result := FindField(reflect.ValueOf(p2), "A")
		require.Equal(t, 1, result.Interface())
	}
	{
		result := FindField(reflect.ValueOf(p2), "B")
		require.False(t, result.IsValid())
	}
}
