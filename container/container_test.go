package container

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestContainer_BindAndResolve(t *testing.T) {
	c := NewContainer()
	type Test struct {
		A int
	}
	testProvider := func() *Test {
		return &Test{
			A: 1,
		}
	}

	t.Run("test nest resolve", func(t *testing.T) {
		type Test2 struct {
			A int
		}
		type Test3 struct {
			A int
		}
		test2Provider := func(test1 *Test) *Test2 {
			return &Test2{
				A: test1.A,
			}
		}
		test3Provider := func(test2 *Test) *Test3 {
			return &Test3{
				A: test2.A,
			}
		}

		require.NoError(t, c.Bind(reflect.TypeOf((*Test)(nil)), reflect.ValueOf(testProvider)))
		require.NoError(t, c.Bind(reflect.TypeOf((*Test2)(nil)), reflect.ValueOf(test2Provider)))
		require.NoError(t, c.Bind(reflect.TypeOf((*Test3)(nil)), reflect.ValueOf(test3Provider)))

		tmp, err := c.Resolve(reflect.TypeOf((*Test3)(nil)))
		require.NoError(t, err)

		test3, ok := tmp.Interface().(*Test3)
		require.True(t, ok)
		require.Equal(t, 1, test3.A)
	})

	t.Run("test struct ptr", func(t *testing.T) {
		require.NoError(t, c.Bind(reflect.TypeOf((*Test)(nil)), reflect.ValueOf(testProvider), NoSingleton()))

		tmp, err := c.Resolve(reflect.TypeOf((*Test)(nil)))
		require.Nil(t, err)

		test, ok := tmp.Interface().(*Test)
		require.True(t, ok)
		require.Equal(t, 1, test.A)

		tmp, err = c.Resolve(reflect.TypeOf((*Test)(nil)))
		require.Nil(t, err)

		test2, ok := tmp.Interface().(*Test)
		require.True(t, ok)
		require.Equal(t, 1, test2.A)

		require.NotEqual(t, fmt.Sprintf("%p", test), fmt.Sprintf("%p", test2))
	})

	t.Run("test struct ptr singleton", func(t *testing.T) {
		require.NoError(t, c.Bind(reflect.TypeOf((*Test)(nil)), reflect.ValueOf(testProvider)))

		tmp, err := c.Resolve(reflect.TypeOf((*Test)(nil)))
		require.NoError(t, err)

		test, ok := tmp.Interface().(*Test)
		require.True(t, ok)

		tmp, err = c.Resolve(reflect.TypeOf((*Test)(nil)))
		require.NoError(t, err)

		test2, ok := tmp.Interface().(*Test)
		require.True(t, ok)

		require.Equal(t, fmt.Sprintf("%p", test), fmt.Sprintf("%p", test2))
	})

	t.Run("test func", func(t *testing.T) {
		var a int
		f := func(a *int) { *a = 1 }
		funcProvider := func() func(*int) {
			return f
		}
		require.NoError(t, c.Bind(reflect.TypeOf(f), reflect.ValueOf(funcProvider)))

		tmp, err := c.Resolve(reflect.TypeOf(f))
		require.NoError(t, err)

		tmp.Interface().(func(*int))(&a)
		require.Equal(t, 1, a)
	})

	t.Run("test alias", func(t *testing.T) {
		type myTest Test
		require.NoError(t, c.Bind(reflect.TypeOf((*myTest)(nil)), reflect.ValueOf(testProvider)))

		tmp, err := c.Resolve(reflect.TypeOf((*myTest)(nil)))
		require.NoError(t, err)

		test, ok := tmp.Interface().(*myTest)
		require.True(t, ok)
		require.Equal(t, 1, test.A)
	})

	t.Run("test bind object", func(t *testing.T) {
		testStruct := &Test{}
		require.NoError(t, c.Bind(reflect.TypeOf((*Test)(nil)), reflect.ValueOf(any(testStruct))))

		tmp, err := c.Resolve(reflect.TypeOf((*Test)(nil)))
		require.Nil(t, err)

		test, ok := tmp.Interface().(*Test)
		require.True(t, ok)

		require.Equal(t, fmt.Sprintf("%p", test), fmt.Sprintf("%p", testStruct))
	})

	t.Run("test bind func", func(t *testing.T) {
		testFunc := func(in string) (out string) {
			return in
		}
		require.NoError(t, c.Bind(reflect.TypeOf(testFunc), reflect.ValueOf(testFunc)))

		tmp, err := c.Resolve(reflect.TypeOf(testFunc))
		require.Nil(t, err)

		test, ok := tmp.Interface().(func(string) string)
		require.True(t, ok)

		require.Equal(t, "1", test("1"))
	})
}
