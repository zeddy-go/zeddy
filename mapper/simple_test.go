package mapper

import (
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestMapper(t *testing.T) {
	t.Run("struct to struct", func(t *testing.T) {
		type a struct {
			A int
			C int
		}

		type c struct {
			C int
		}
		type b struct {
			B *int
			c
		}

		aa := &a{
			A: 1,
		}

		var bb b

		err := NewSimple(WithFieldSimple(func(srcField string) (dstField string) {
			if srcField == "A" {
				return "B"
			}
			return srcField
		})).Map(&bb, aa)
		require.NoError(t, err)
		require.Equal(t, aa.A, *bb.B)
		require.Equal(t, aa.C, bb.C)
	})

	t.Run("struct to struct", func(t *testing.T) {
		type c struct {
			C int
		}

		type a struct {
			A int
			c
		}

		type b struct {
			B *int
			c
		}

		aa := &a{
			A: 1,
			c: c{
				C: 1,
			},
		}

		var bb b

		err := NewSimple(WithFieldSimple(func(srcField string) (dstField string) {
			if srcField == "A" {
				return "B"
			}
			return srcField
		})).Map(&bb, aa)
		require.NoError(t, err)
		require.Equal(t, aa.A, *bb.B)
		require.Equal(t, aa.C, bb.C)
	})

	t.Run("struct to map", func(t *testing.T) {
		type a struct {
			A int
		}

		aa := a{
			A: 1,
		}

		var bb map[string]any

		err := NewSimple(WithFieldSimple(func(srcField string) (dstField string) {
			return strings.ToLower(srcField)
		})).Map(&bb, aa)
		require.NoError(t, err)
		require.Equal(t, aa.A, bb["a"])
	})

	t.Run("map to struct", func(t *testing.T) {
		type a struct {
			A int
		}

		var aa a

		bb := map[string]any{
			"a": "1",
		}

		err := NewSimple(WithFieldSimple(func(srcField string) (dstField string) {
			return strings.ToUpper(srcField)
		})).Map(&aa, bb)
		require.NoError(t, err)
		require.Equal(t, 1, aa.A)
	})

	t.Run("map slice", func(t *testing.T) {
		type a struct {
			A int
		}

		type b struct {
			A int
		}

		aa := []a{
			{
				A: 1,
			},
		}

		var bb []b

		err := NewSimple().MapSlice(&bb, aa)
		require.NoError(t, err)
		require.Equal(t, aa[0].A, bb[0].A)
	})
}

type s1 struct {
	A int
	B string
	C bool
}

type s2 struct {
	A int
	B string
	C bool
}

func BenchmarkCopier(b *testing.B) {
	struct1 := s1{
		A: 1,
		B: "2",
		C: true,
	}
	for i := 0; i < b.N; i++ {
		var struct2 s2
		err := copier.Copy(&struct2, struct1)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkSimpleMap(b *testing.B) {
	struct1 := s1{
		A: 1,
		B: "2",
		C: true,
	}
	for i := 0; i < b.N; i++ {
		var struct2 s2
		err := SimpleMap(&struct2, struct1)
		if err != nil {
			b.Error(err)
		}
	}
}
