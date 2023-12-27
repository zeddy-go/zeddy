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

		c.Register(NewStuff(testProvider))
		c.Register(NewStuff(test2Provider))
		c.Register(NewStuff(test3Provider))

		tmp, err := c.Resolve(reflect.TypeOf(&Test3{}))
		require.Nil(t, err)

		test3, ok := tmp.(*Test3)
		require.True(t, ok)
		require.Equal(t, 1, test3.A)
	})

	t.Run("test struct ptr", func(t *testing.T) {
		t.Parallel()
		stuff := NewStuff(testProvider)
		c.Register(stuff)

		tmp, err := c.Resolve(reflect.TypeOf(&Test{}))
		require.Nil(t, err)

		test, ok := tmp.(*Test)
		require.True(t, ok)
		require.Equal(t, 1, test.A)

		tmp, err = c.Resolve(reflect.TypeOf(&Test{}))
		require.Nil(t, err)

		test2, ok := tmp.(*Test)
		require.True(t, ok)
		require.Equal(t, 1, test2.A)

		require.NotEqual(t, fmt.Sprintf("%p", test), fmt.Sprintf("%p", test2))
	})

	t.Run("test struct ptr singleton", func(t *testing.T) {
		t.Parallel()
		stuff := NewStuff(testProvider, WithSingleton(), WithKey("test1"))
		c.Register(stuff)

		tmp, err := c.Resolve(reflect.TypeOf(&Test{}), "test1")
		require.Nil(t, err)

		test, ok := tmp.(*Test)
		require.True(t, ok)

		tmp, err = c.Resolve(reflect.TypeOf(&Test{}), "test1")
		require.Nil(t, err)

		test2, ok := tmp.(*Test)
		require.True(t, ok)

		require.Equal(t, fmt.Sprintf("%p", test), fmt.Sprintf("%p", test2))
	})

	t.Run("test func", func(t *testing.T) {
		t.Parallel()
		var a int
		f := func(a *int) { *a = 1 }
		funcProvider := func() func(*int) {
			return f
		}
		stuff := NewStuff(funcProvider)
		c.Register(stuff)

		tmp, err := c.Resolve(reflect.TypeOf(f))
		require.Nil(t, err)

		tmp.(func(*int))(&a)
		require.Equal(t, 1, a)
	})
}

func TestContainer_PkgBindAndResolve(t *testing.T) {
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

		Register(testProvider)
		Register(test2Provider)
		Register(test3Provider)

		test3, err := Resolve[*Test3]()
		require.Nil(t, err)
		require.Equal(t, 1, test3.A)
	})

	t.Run("test struct ptr", func(t *testing.T) {
		t.Parallel()
		Register(testProvider)

		test, err := Resolve[*Test]()
		require.Nil(t, err)
		require.Equal(t, 1, test.A)

		test2, err := Resolve[*Test]()
		require.Nil(t, err)
		require.Equal(t, 1, test2.A)

		require.NotEqual(t, fmt.Sprintf("%p", test), fmt.Sprintf("%p", test2))
	})

	t.Run("test struct ptr singleton", func(t *testing.T) {
		t.Parallel()
		Register(testProvider, WithSingleton(), WithKey("test1"))

		test, err := Resolve[*Test]("test1")
		require.Nil(t, err)

		test2, err := Resolve[*Test]("test1")
		require.Nil(t, err)

		require.Equal(t, fmt.Sprintf("%p", test), fmt.Sprintf("%p", test2))
	})

	t.Run("test func", func(t *testing.T) {
		t.Parallel()
		var a int
		f := func(a *int) { *a = 1 }
		funcProvider := func() func(*int) {
			return f
		}
		Register(funcProvider)

		tmp, err := Resolve[func(a *int)]()
		require.Nil(t, err)

		tmp(&a)
		require.Equal(t, 1, a)
	})
}

type testStruct struct {
}

type testStruct2 testStruct

func TestXxx(t *testing.T) {
	a := testStruct{}
	_ = reflect.ValueOf(a).Convert(reflect.TypeOf((*testStruct2)(nil)).Elem()).Interface().(testStruct2)
}
