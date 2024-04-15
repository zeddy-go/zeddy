package mapper

import (
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestFindField(t *testing.T) {
	t.Run("case", func(t *testing.T) {
		type s1 struct {
			AA bool
		}
		type s2 struct {
			Aa bool
		}
		ss1 := s1{}
		ss2 := s2{}
		vS1 := reflect.ValueOf(&ss1).Elem()
		vS2 := reflect.ValueOf(ss2)
		f := findField(vS1, vS2.Type().Field(0), vS2.Field(0), true)
		require.False(t, f.IsValid())

		f = findField(vS1, vS2.Type().Field(0), vS2.Field(0), false)
		f.Set(reflect.ValueOf(true))
		require.Equal(t, true, ss1.AA)
	})
	t.Run("anonymous", func(t *testing.T) {
		type a struct {
			C bool
		}
		type s1 struct {
			a
		}
		type s2 struct {
			a
		}
		ss1 := s1{}
		ss2 := s2{}
		vS1 := reflect.ValueOf(&ss1).Elem()
		vS2 := reflect.ValueOf(ss2)
		f := findField(vS1, vS2.Type().Field(0), vS2.Field(0), false)
		require.Equal(t, reflect.Struct, f.Kind())

		f = findField(f, vS2.Field(0).Type().Field(0), vS2.Field(0).Field(0), false)
		f.Set(reflect.ValueOf(true))
		require.Equal(t, true, ss1.a.C)
	})
	t.Run("anonymous_pointer", func(t *testing.T) {
		type A struct {
			C bool
		}
		type s1 struct {
			*A
		}
		type s2 struct {
			*A
		}
		ss1 := s1{}
		ss2 := s2{
			A: &A{},
		}
		vS1 := reflect.ValueOf(&ss1).Elem()
		vS2 := reflect.ValueOf(ss2)
		f := findField(vS1, vS2.Type().Field(0), vS2.Field(0), false)
		require.Equal(t, reflect.Pointer, f.Kind())

		f = reflect.New(f.Type().Elem())
		vS1.Field(0).Set(f)
		f = findField(f.Elem(), vS2.Field(0).Elem().Type().Field(0), vS2.Field(0).Elem().Field(0), false)
		f.Set(reflect.ValueOf(true))
		require.Equal(t, true, ss1.A.C)
	})
}

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
