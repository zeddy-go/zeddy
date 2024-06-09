package errx

import (
	"context"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestXxx(t *testing.T) {
	prefix, _ := filepath.Abs(".")
	prefix = strings.Replace(prefix, "\\", "/", -1)
	err := New("123")
	err2 := WrapWithSkip(err, "321", 0)
	require.Equal(t, prefix+"/err_test.go:22:321\n"+prefix+"/err_test.go:21:123\n\n", fmt.Sprintf("%#v\n", err2))
	require.Equal(t, "321: 123", err2.Error())
	require.Equal(t, prefix+"/err_test.go:22:321\n\n", fmt.Sprintf("%+v\n", err2))
}

func TestGrpcStatus(t *testing.T) {
	e := Wrap(Wrap(errors.New("test"), "test1"), "test2").(*Errx)
	e2 := NewFromStatus(e.GRPCStatus())
	require.NotNil(t, e2)
	prefix, _ := filepath.Abs(".")
	prefix = strings.Replace(prefix, "\\", "/", -1)
	expect := prefix + "/err_test.go:29:test2\n" + prefix + "/err_test.go:29:test1\n" + "test"
	require.Equal(t, expect, fmt.Sprintf("%#v", e2))
}

type testServer struct {
	UnimplementedTestServer
}

func (t testServer) Test(context.Context, *NoContent) (*NoContent, error) {
	return nil, New("test")
}

func TestGrpcTrans(t *testing.T) {
	lis, err := net.Listen("tcp", "0.0.0.0:8888")
	require.NoError(t, err)
	s := grpc.NewServer()
	RegisterTestServer(s, testServer{})
	go s.Serve(lis)

	time.Sleep(time.Second)

	cc, err := grpc.Dial("0.0.0.0:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	c := NewTestClient(cc)
	_, err = c.Test(context.Background(), &NoContent{})
	fmt.Printf("%+v\n", err)
}

func TestErrorsIs(t *testing.T) {
	e1 := errors.New("test")
	e2 := fmt.Errorf("test2: %w", e1)
	require.True(t, Is(e2, e1))
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
		Detail: "test",
	})

	require.Equal(t, "test", e.MustGet(Detail))
}
