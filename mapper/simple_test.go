package mapper

import (
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestSimpleMap(t *testing.T) {
	type CommonField struct {
		ID uint64
	}

	type UserEntity struct {
		Id       uint64
		Username string
	}

	type UserModel struct {
		CommonField
		UserName string
	}

	t.Run("entity_to_model", func(t *testing.T) {
		e := UserEntity{
			Id:       1,
			Username: "zed",
		}
		var m UserModel
		err := SimpleMap(&m, e)
		require.NoError(t, err)
		require.Equal(t, e.Username, m.UserName)
		require.Equal(t, e.Id, m.ID)
	})

	t.Run("model_to_entity", func(t *testing.T) {
		m := UserModel{
			CommonField: CommonField{
				ID: 1,
			},
			UserName: "zed",
		}
		var e UserEntity
		err := SimpleMap(&e, m)
		require.NoError(t, err)
		require.Equal(t, m.UserName, e.Username)
		require.Equal(t, m.ID, e.Id)
	})

	t.Run("anonymous", func(t *testing.T) {
		type CommonField struct {
			ID uint64
		}

		type UserEntity struct {
			CommonField
			Username string
		}

		type UserModel struct {
			CommonField
			UserName string
		}

		m := UserModel{
			CommonField: CommonField{
				ID: 1,
			},
			UserName: "zed",
		}
		var e UserEntity
		err := SimpleMap(&e, m)
		require.NoError(t, err)
		require.Equal(t, m.UserName, e.Username)
		require.Equal(t, m.ID, e.ID)
	})
}

func TestFieldInValid(t *testing.T) {
	type CommonField struct {
		C bool
	}

	type s1 struct {
		A int
		B string
		CommonField
	}

	type s2 struct {
		A int
		B string
		CommonField
	}

	struct1 := s1{
		A:           1,
		B:           "2",
		CommonField: CommonField{C: true},
	}
	var struct2 s2
	err := SimpleMap(&struct2, struct1)
	if err != nil {
		require.NoError(t, err)
	}
	require.Equal(t, struct1.A, struct2.A)
	require.Equal(t, struct1.B, struct2.B)
	require.Equal(t, struct1.C, struct2.C)
}

func TestSimpleMap2(t *testing.T) {
	type m struct {
		ID       uint64
		Username string
		Password string
	}
	type e struct {
		ID       uint64
		Username string
		Password string
	}

	model := &m{
		ID:       1,
		Username: "zed",
		Password: "zed",
	}
	var entity e
	err := SimpleMap(&entity, &model)
	require.NoError(t, err)
	require.Equal(t, model.ID, entity.ID)
	require.Equal(t, model.Username, entity.Username)
	require.Equal(t, model.Password, entity.Password)
}

func TestSimpleMap3(t *testing.T) {
	type Common1 struct {
		C bool
	}
	type Common2 struct {
		C bool
	}
	type s1 struct {
		A      int
		Common *Common1
	}
	type s2 struct {
		A      int
		Common Common2
	}

	struct1 := s1{
		A:      1,
		Common: &Common1{C: true},
	}
	var struct2 s2
	err := SimpleMap(&struct2, struct1)
	require.NoError(t, err)
	require.Equal(t, struct1.A, struct2.A)
	require.Equal(t, struct1.Common.C, struct2.Common.C)
}

func TestSimpleMap4(t *testing.T) {
	type S1 struct {
		A string
	}
	type S2 struct {
		A []string
	}
	s1 := S1{}
	var s2 S2
	err := SimpleMap(&s2, s1)
	require.NoError(t, err)
	require.Equal(t, 0, len(s2.A))
}

func TestSimpleMap5(t *testing.T) {
	type Common1 struct {
		C bool
	}
	type s1 struct {
		A int
		*Common1
	}
	type s2 struct {
		A int
		*Common1
	}

	struct1 := s1{
		A:       1,
		Common1: &Common1{C: true},
	}
	var struct2 s2
	err := SimpleMap(&struct2, struct1)
	require.NoError(t, err)
	require.Equal(t, struct1.A, struct2.A)
	require.Equal(t, struct1.Common1.C, struct2.Common1.C)
}

func TestSimpleMap6(t *testing.T) {
	type Common1 struct {
		C bool
	}
	type Common2 struct {
		C bool
	}
	type s1 struct {
		A      int
		Common *Common1
	}
	type s2 struct {
		A      int
		Common *Common2
	}

	struct1 := s1{
		A:      1,
		Common: &Common1{C: true},
	}
	var struct2 s2
	err := SimpleMap(&struct2, struct1)
	require.NoError(t, err)
	require.Equal(t, struct1.A, struct2.A)
	require.Equal(t, struct1.Common.C, struct2.Common.C)
}

func TestSimpleMapSliceValue(t *testing.T) {
	type s1 struct {
		A int
	}
	type s2 struct {
		A int
	}

	slice1 := []s1{
		{A: 1},
	}
	slice2 := make([]s2, 0)
	err := SimpleMapSliceValue(reflect.ValueOf(&slice2), reflect.ValueOf(slice1))
	require.NoError(t, err)
	require.Equal(t, slice1[0].A, slice2[0].A)
}

func TestSimpleMapSliceValue2(t *testing.T) {
	type s1 struct {
		A int
	}
	type s2 struct {
		A int
	}

	slice1 := []s1{
		{A: 1},
	}
	var slice2 []s2
	err := SimpleMapSliceValue(reflect.ValueOf(&slice2), reflect.ValueOf(slice1))
	require.NoError(t, err)
	require.Equal(t, slice1[0].A, slice2[0].A)
}

func TestSimpleMapSlice(t *testing.T) {
	type struct1 struct {
		A bool
	}
	type struct2 struct {
		A bool
	}
	s1 := []struct1{{A: true}}
	var s2 []struct2
	err := SimpleMapSlice(&s2, s1)
	require.NoError(t, err)
	require.Equal(t, s1[0].A, s2[0].A)
}

func TestSimpleMap7(t *testing.T) {
	type c1 struct {
		B int
	}
	type c2 struct {
		B int
	}
	type s1 struct {
		Bs []*c1
	}
	type s2 struct {
		Bs []*c2
	}

	struct1 := s1{
		Bs: []*c1{{B: 1}},
	}
	var struct2 s2
	err := SimpleMap(&struct2, struct1)
	require.NoError(t, err)
	require.Equal(t, struct1.Bs[0].B, struct2.Bs[0].B)
}

func TestSimpleMapSlice3(t *testing.T) {
	type Common2 struct {
		ID uint64 `json:"id,string"`
	}

	type struct2 struct {
		Types         []int    `json:"types"`
		CurrentNumber *Common2 `json:"currentNumber"`
	}

	type Common1 struct {
		ID uint64
	}

	type struct1 struct {
		Types         []int
		CurrentNumber *Common1
	}

	s1 := []struct1{{Types: []int{}, CurrentNumber: &Common1{ID: 1}}}
	var s2 []*struct2
	err := SimpleMapSlice(&s2, s1)
	require.NoError(t, err)
	require.Equal(t, uint64(1), s2[0].CurrentNumber.ID)
}

func BenchmarkCopier(b *testing.B) {
	type CommonField struct {
		C bool
	}

	type s1 struct {
		A int
		B string
		CommonField
	}

	type s2 struct {
		A int
		B string
		CommonField
	}

	b.Run("copier", func(b *testing.B) {
		struct1 := s1{
			A:           1,
			B:           "2",
			CommonField: CommonField{C: true},
		}
		for i := 0; i < b.N; i++ {
			var struct2 s2
			err := copier.Copy(&struct2, struct1)
			if err != nil {
				require.NoError(b, err)
			}
			require.Equal(b, struct1.A, struct2.A)
			require.Equal(b, struct1.B, struct2.B)
			require.Equal(b, struct1.C, struct2.C)
		}
	})

	b.Run("SimpleMap", func(b *testing.B) {
		struct1 := s1{
			A:           1,
			B:           "2",
			CommonField: CommonField{C: true},
		}
		for i := 0; i < b.N; i++ {
			var struct2 s2
			err := SimpleMap(&struct2, struct1)
			if err != nil {
				require.NoError(b, err)
			}
			require.Equal(b, struct1.A, struct2.A)
			require.Equal(b, struct1.B, struct2.B)
			require.Equal(b, struct1.C, struct2.C)
		}
	})
}
