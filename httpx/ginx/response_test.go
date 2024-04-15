package ginx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/require"
	"github.com/zeddy-go/zeddy/errx"
	"net/http"
	"net/http/httptest"
	"testing"
)

func do(resp IResponse[*gin.Context]) (response *httptest.ResponseRecorder) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		resp.Do(c)
	})

	response = httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(response, request)
	return
}

func TestJsonResponse(t *testing.T) {
	r := do(NewRestfulResponse())
	require.Equal(t, 204, r.Code)

	r = do(NewRestfulResponse().SetData(gin.H{"test": true}))
	require.Equal(t, 200, r.Code)
	require.Equal(t, `{"data":{"test":true},"message":""}`, r.Body.String())

	r = do(NewRestfulResponse().SetData(gin.H{"test": true}).SetMeta(&Meta{Total: 1}))
	require.Equal(t, 200, r.Code)
	require.Equal(t, `{"data":{"test":true},"message":"","meta":{"total":1}}`, r.Body.String())

	r = do(NewRestfulResponse().SetData(gin.H{"test": true}).SetMeta(&Meta{Total: 1}).SetError(errors.New("test")))
	require.Equal(t, http.StatusInternalServerError, r.Code)
	require.Equal(t, `{"data":null,"message":"test"}`, r.Body.String())

	r = do(NewRestfulResponse().SetError(errors.New("test")))
	require.Equal(t, http.StatusInternalServerError, r.Code)
	require.Equal(t, `{"data":null,"message":"test"}`, r.Body.String())

	r = do(NewRestfulResponse().SetError(errx.New("test", errx.WithCode(5000))))
	require.Equal(t, http.StatusInternalServerError, r.Code)
	require.Equal(t, `{"data":null,"message":"test"}`, r.Body.String())

	r = do(NewRestfulResponse().SetError(errx.New("test", errx.WithCode(400))))
	require.Equal(t, http.StatusBadRequest, r.Code)
	require.Equal(t, `{"data":null,"message":"test"}`, r.Body.String())

	type a struct {
		A int `json:"a" binding:"required"`
		B int `json:"b" binding:"required"`
	}
	r = do(NewRestfulResponse().SetError(binding.Validator.ValidateStruct(&a{})))
	require.Equal(t, http.StatusUnprocessableEntity, r.Code)
	require.Equal(t, `{"data":null,"message":"Key: 'a.A' Error:Field validation for 'A' failed on the 'required' tag\nKey: 'a.B' Error:Field validation for 'B' failed on the 'required' tag"}`, r.Body.String())
}

func TestFileResponse(t *testing.T) {
	content := []byte{1, 2, 3, 4}
	r := do(NewRestfulResponse().SetFile(NewFile(content)))
	require.Equal(t, http.StatusOK, r.Code)
	require.Equal(t, content, r.Body.Bytes())
	require.Equal(t, "application/octet-stream", r.Header().Get("Content-Type"))
	require.Equal(t, "", r.Header().Get("Content-Disposition"))

	r = do(NewRestfulResponse().SetFile(NewFile(content, WithFileName("test.jpg"), WithMimeType("images/jpeg"))))
	require.Equal(t, http.StatusOK, r.Code)
	require.Equal(t, content, r.Body.Bytes())
	require.Equal(t, "images/jpeg", r.Header().Get("Content-Type"))
	require.Equal(t, "attachment; filename=\"test.jpg\"", r.Header().Get("Content-Disposition"))
}
