package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func a(m IMeta) IMeta {
	return nil
}

func TestXxx(t *testing.T) {

	mt := reflect.TypeOf(a)

	println(mt.Out(0).Name())
}

func TestContext(t *testing.T) {
	c := &gin.Context{
		Keys: map[string]any{
			"test": true,
		},
	}

	require.Equal(t, reflect.ValueOf(c).Type(), reflect.TypeOf(&gin.Context{}))
}

func TestParse(t *testing.T) {
	type req struct {
		ID   string `uri:"id" binding:"required"`
		Name string `form:"name" binding:"required"`
		Num  int    `json:"num" binding:"required"`
	}
	r := gin.Default()
	r.POST("/test/:id", func(c *gin.Context) {
		result, err := parseParam(c, reflect.TypeOf(&req{}))
		require.NoError(t, err)
		request := result.Interface().(*req)
		require.Equal(t, "1", request.ID)
		require.Equal(t, "2", request.Name)
		require.Equal(t, 4, request.Num)
	})

	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/test/1?name=2", strings.NewReader(`{"num": 4}`))
	request.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, request)

	require.Equal(t, 200, w.Code)
}

func TestParse2(t *testing.T) {
	type req struct {
		ID    string `uri:"id" binding:"required"`
		Name  string `form:"name" binding:"required"`
		Name2 string `form:"name2" binding:"required"`
	}
	r := gin.Default()
	r.POST("/test/:id", func(c *gin.Context) {
		result, err := parseParam(c, reflect.TypeOf(&req{}))
		require.NoError(t, err)
		request := result.Interface().(*req)
		require.Equal(t, "1", request.ID)
		require.Equal(t, "2", request.Name)
		require.Equal(t, "3", request.Name2)
	})

	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/test/1?name=2", strings.NewReader(`name2=3`))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, request)

	require.Equal(t, 200, w.Code)
}

func TestReflectOfParams(t *testing.T) {
	a := func(a int, b string) {}

	fmt.Printf("%+v\n", reflect.TypeOf(a).In(1).Name())
}

func TestResponse(t *testing.T) {
	r := response().(*JsonResponse)
	require.Nil(t, r.err)
	require.Nil(t, r.response)

	r = response(reflect.ValueOf(errors.New("test"))).(*JsonResponse)
	require.NotNil(t, r.err)
	require.Nil(t, r.response)

	r = response(reflect.ValueOf(gin.H{"test": true})).(*JsonResponse)
	require.Nil(t, r.err)
	require.NotNil(t, r.response)
	require.NotNil(t, r.response.Data)

	r = response(reflect.ValueOf(gin.H{"test": true}), reflect.ValueOf(nil)).(*JsonResponse)
	require.Nil(t, r.err)
	require.NotNil(t, r.response)
	require.NotNil(t, r.response.Data)

	r = response(reflect.ValueOf(gin.H{"test": true}), reflect.ValueOf(errors.New("test"))).(*JsonResponse)
	require.NotNil(t, r.err)
	require.Nil(t, r.response)

	r = response(reflect.ValueOf(&Meta{Total: 1}), reflect.ValueOf(gin.H{"test": true}), reflect.ValueOf(errors.New("test"))).(*JsonResponse)
	require.NotNil(t, r.err)
	require.Nil(t, r.response)

	r = response(reflect.ValueOf(&Meta{Total: 1}), reflect.ValueOf(gin.H{"test": true}), reflect.ValueOf(nil)).(*JsonResponse)
	require.Nil(t, r.err)
	require.NotNil(t, r.response)
	require.NotNil(t, r.response.Data)
	require.NotNil(t, r.response.Meta)
}
