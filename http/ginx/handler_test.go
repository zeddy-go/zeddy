package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

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

func TestResponse(t *testing.T) {
	r := parseAndResponse().(*JsonResponse)
	require.Nil(t, r.Err)
	require.Nil(t, r.Response)

	r = parseAndResponse(reflect.ValueOf(errors.New("test"))).(*JsonResponse)
	require.NotNil(t, r.Err)
	require.Nil(t, r.Response)

	r = parseAndResponse(reflect.ValueOf(gin.H{"test": true})).(*JsonResponse)
	require.Nil(t, r.Err)
	require.NotNil(t, r.Response)
	require.NotNil(t, r.Response.Data)

	r = parseAndResponse(reflect.ValueOf(gin.H{"test": true}), reflect.ValueOf(nil)).(*JsonResponse)
	require.Nil(t, r.Err)
	require.NotNil(t, r.Response)
	require.NotNil(t, r.Response.Data)

	r = parseAndResponse(reflect.ValueOf(gin.H{"test": true}), reflect.ValueOf(errors.New("test"))).(*JsonResponse)
	require.NotNil(t, r.Err)
	require.Nil(t, r.Response)

	r = parseAndResponse(reflect.ValueOf(&Meta{Total: 1}), reflect.ValueOf(gin.H{"test": true}), reflect.ValueOf(errors.New("test"))).(*JsonResponse)
	require.NotNil(t, r.Err)
	require.Nil(t, r.Response)

	r = parseAndResponse(reflect.ValueOf(&Meta{Total: 1}), reflect.ValueOf(gin.H{"test": true}), reflect.ValueOf(nil)).(*JsonResponse)
	require.Nil(t, r.Err)
	require.NotNil(t, r.Response)
	require.NotNil(t, r.Response.Data)
	require.NotNil(t, r.Response.Meta)
}
