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

	t.Run("test bind instance", func(t *testing.T) {
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

	t.Run("test bind value", func(t *testing.T) {
		testStruct := Test{}
		err := c.Bind(reflect.TypeOf((*Test)(nil)), reflect.ValueOf(any(testStruct)))
		require.Error(t, err)
		require.Contains(t, err.Error(), "can not bind <container.Test> to <*container.Test>")
	})

	t.Run("test bind value no singleton", func(t *testing.T) {
		testStruct := &Test{}
		err := c.Bind(reflect.TypeOf((*Test)(nil)), reflect.ValueOf(testStruct), NoSingleton())
		require.NoError(t, err)

		tmp, err := c.Resolve(reflect.TypeOf((*Test)(nil)))
		require.Nil(t, err)

		test, ok := tmp.Interface().(*Test)
		require.True(t, ok)

		require.NotEqual(t, fmt.Sprintf("%p", test), fmt.Sprintf("%p", testStruct))

		testStruct = test

		tmp, err = c.Resolve(reflect.TypeOf((*Test)(nil)))
		require.Nil(t, err)

		test, ok = tmp.Interface().(*Test)
		require.True(t, ok)

		require.NotEqual(t, fmt.Sprintf("%p", test), fmt.Sprintf("%p", testStruct))
	})

	t.Run("test bind error1", func(t *testing.T) {
		err := c.Bind(reflect.TypeOf((*Test)(nil)), reflect.ValueOf(func() {}))
		require.Error(t, err)
		require.Contains(t, err.Error(), "can not bind <func()> to <*container.Test>")
	})

	t.Run("test bind convertable", func(t *testing.T) {
		type myTest Test
		test := &Test{}
		err := c.Bind(reflect.TypeOf((*myTest)(nil)), reflect.ValueOf(test))
		require.NoError(t, err)

		tmp, err := c.Resolve(reflect.TypeOf((*myTest)(nil)))
		require.Nil(t, err)

		testStruct, ok := tmp.Interface().(*myTest)
		require.True(t, ok)

		require.Equal(t, fmt.Sprintf("%p", test), fmt.Sprintf("%p", testStruct))
	})

	t.Run("test bind convertable not singleton", func(t *testing.T) {
		type myTest Test
		test := &Test{}
		err := c.Bind(reflect.TypeOf((*myTest)(nil)), reflect.ValueOf(test), NoSingleton())
		require.NoError(t, err)

		tmp, err := c.Resolve(reflect.TypeOf((*myTest)(nil)))
		require.Nil(t, err)

		testStruct, ok := tmp.Interface().(*myTest)
		require.True(t, ok)

		require.NotEqual(t, fmt.Sprintf("%p", test), fmt.Sprintf("%p", testStruct))

		tmp, err = c.Resolve(reflect.TypeOf((*myTest)(nil)))
		require.Nil(t, err)

		testStruct2, ok := tmp.Interface().(*myTest)
		require.True(t, ok)

		require.NotEqual(t, fmt.Sprintf("%p", testStruct), fmt.Sprintf("%p", testStruct2))
	})

	t.Run("cycleResolve1", func(t *testing.T) {
		err := Bind[*Struct1](NewStruct1)
		require.NoError(t, err)
		err = Bind[*Struct2](NewStruct2)
		require.NoError(t, err)

		s2, err := Resolve[*Struct2]()
		require.NoError(t, err)
		require.Equal(t, s2.A, 2)
		require.Equal(t, s2.S1.A, 1)

		s1, err := Resolve[*Struct1]()
		require.NoError(t, err)
		require.Equal(t, s1.A, 1)
		require.Equal(t, s1.S2.A, 2)

		require.Same(t, s1.S2, s2)
		require.Same(t, s1, s2.S1)
	})

	t.Run("cycleResolve2", func(t *testing.T) {
		err := Bind[*Struct11](NewStruct11)
		require.NoError(t, err)
		err = Bind[*Struct12](NewStruct12)
		require.NoError(t, err)
		err = Bind[*Struct13](NewStruct13)
		require.NoError(t, err)

		s11, err := Resolve[*Struct11]()
		require.NoError(t, err)

		s12, err := Resolve[*Struct12]()
		require.NoError(t, err)

		s13, err := Resolve[*Struct13]()
		require.NoError(t, err)

		require.Same(t, s11, s13.s11)
		require.Same(t, s12, s11.s12)
		require.Same(t, s13, s12.s13)
	})

	//t.Run("cycleResolve3", func(t *testing.T) {
	//	err := Bind[*Struct111](NewStruct111)
	//	require.NoError(t, err)
	//	err = Bind[*Struct112](NewStruct112)
	//	require.NoError(t, err)
	//	err = Bind[*Struct113](NewStruct113)
	//	require.NoError(t, err)
	//
	//	s111, err := Resolve[*Struct111]()
	//	require.NoError(t, err)
	//
	//	s112, err := Resolve[*Struct112]()
	//	require.NoError(t, err)
	//
	//	s113, err := Resolve[*Struct113]()
	//	require.NoError(t, err)
	//
	//	require.Same(t, s111, s113.s111)
	//	require.Same(t, s111, s112.s111)
	//	require.Same(t, s112, s111.s112)
	//	require.Same(t, s112, s113.s112)
	//	require.Same(t, s113, s112.s113)
	//	require.Same(t, s113, s111.s113)
	//})
}

func NewStruct2(s1 *Struct1) *Struct2 {
	return &Struct2{
		S1: s1,
		A:  2,
	}
}

type Struct2 struct {
	S1 *Struct1
	A  int
}

func NewStruct1(s2 *Struct2) *Struct1 {
	return &Struct1{
		S2: s2,
		A:  1,
	}
}

type Struct1 struct {
	S2 *Struct2
	A  int
}

func NewStruct11(s12 *Struct12) *Struct11 {
	return &Struct11{
		s12: s12,
	}
}

type Struct11 struct {
	s12 *Struct12
}

func NewStruct12(s13 *Struct13) *Struct12 {
	return &Struct12{
		s13: s13,
	}
}

type Struct12 struct {
	s13 *Struct13
}

func NewStruct13(s11 *Struct11) *Struct13 {
	return &Struct13{
		s11: s11,
	}
}

func NewStruct111(s112 *Struct112, s113 *Struct113) *Struct111 {
	return &Struct111{
		s112: s112,
		s113: s113,
	}
}

type Struct111 struct {
	s112 *Struct112
	s113 *Struct113
}

func NewStruct112(s111 *Struct111, s113 *Struct113) *Struct112 {
	return &Struct112{
		s111: s111,
		s113: s113,
	}
}

type Struct112 struct {
	s111 *Struct111
	s113 *Struct113
}

func NewStruct113(s112 *Struct112, s111 *Struct111) *Struct113 {
	return &Struct113{
		s112: s112,
		s111: s111,
	}
}

type Struct113 struct {
	s111 *Struct111
	s112 *Struct112
}

type Struct13 struct {
	s11 *Struct11
}

func testFunc() int {
	return 1
}

func TestXxx(t *testing.T) {
	s1 := reflect.ValueOf(NewStruct1).Call([]reflect.Value{reflect.ValueOf((*Struct2)(nil))})[0]
	s2 := reflect.ValueOf(NewStruct2).Call([]reflect.Value{s1})[0]
	ss1 := s1
	s1.Elem().Set(reflect.ValueOf(NewStruct1).Call([]reflect.Value{s2})[0].Elem())

	println(s2.Interface().(*Struct2).S1.A)
	println(s1.Interface().(*Struct1).S2.A)
	fmt.Printf("%p %p\n", ss1.Interface().(*Struct1), s1.Interface().(*Struct1))
}
