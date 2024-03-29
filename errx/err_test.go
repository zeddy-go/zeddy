package errx

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGrpcStatus(t *testing.T) {
	e := Wrap(errors.New("test"), "test1").(*Errx)
	e2 := NewFromStatus(e.GRPCStatus())
	fmt.Printf("%+v\n", e2 == nil)
	fmt.Printf("%#v\n", e2)
}

func TestErrorsIs(t *testing.T) {
	e1 := errors.New("test")
	e2 := fmt.Errorf("test2: %w", e1)
	require.True(t, Is(e2, e1))
}

func TestXxx(t *testing.T) {
	err := New("123")
	err2 := WrapWithSkip(err, "321", 0)
	require.Equal(t, "F:/projects/zeddy/zeddy/errx/err_test.go:12:321\nF:/projects/zeddy/zeddy/errx/err_test.go:11:123\n\n", fmt.Sprintf("%#v\n", err2))
	require.Equal(t, "321: 123", err2.Error())
	require.Equal(t, "F:/projects/zeddy/zeddy/errx/err_test.go:12:321\n\n", fmt.Sprintf("%+v\n", err2))
}

func TestErrxIs(t *testing.T) {
	require.True(t, Is(Wrap(New("test"), "test2"), New("test")))
	require.True(t, Is(Wrap(New("test", WithCode(1)), "test2"), New("test3", WithCode(1))))
	require.False(t, Is(Wrap(New("test", WithCode(1)), "test2"), New("test3")))
}

func TestErrorIs(t *testing.T) {
	e := New("test")
	ee := Wrap(e, "test2")
	a := ee
	require.True(t, errors.Is(a, e))
}

func TestChange(t *testing.T) {
	e := New("test").(*Errx)
	e.Set(map[InfoKey]any{
		"detail": "test",
	})

	require.Equal(t, "test", e.MustGet(Detail))
}
