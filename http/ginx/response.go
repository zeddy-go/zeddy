package ginx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zeddy-go/zeddy/errx"
	"gorm.io/gorm"
	"net/http"
)

type IMeta interface {
	GetMeta() any
}

type Meta struct {
	CurrentPage uint `json:"currentPage,omitempty"`
	Total       uint `json:"total,omitempty"`
	LastPage    uint `json:"lastPage,omitempty"`
	PerPage     uint `json:"perPage,omitempty"`
}

func (m *Meta) GetMeta() any {
	return m
}

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Meta    any         `json:"meta,omitempty"`
}

type ResponseWithCode struct {
	Code
	Response
}

type Code struct {
	Code int `json:"code"`
}

// Error 返回错误响应
func Error(c *gin.Context, err error, status int) {
	code, res := parseError(status, err)
	Json(c, res, code, false)
}

func ErrorAbort(c *gin.Context, err error, status int) {
	code, res := parseError(status, err)
	Json(c, res, code, true)
}

func parseError(status int, err error) (code int, res *Response) {
	res = &Response{Message: err.Error()}

	if status != 0 {
		code = status
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		} else if er, ok := err.(errx.Errx); ok && errx.GetErrxField[int](er, errx.Code) != 0 {
			code = errx.GetErrxField[int](er, errx.Code)
			res.Data = errx.GetErrxField[any](er, errx.Detail)
		} else {
			code = http.StatusInternalServerError
		}
	}

	return
}

func Success(c *gin.Context, data any, status int) {
	var (
		code     int
		response *Response
	)

	if status != 0 {
		code = status
	} else {
		if data == nil {
			code = http.StatusNoContent
		} else {
			code = http.StatusOK
		}
	}

	response = &Response{
		Data: data,
	}

	Json(c, response, code, false)
}

func Pagination(c *gin.Context, data interface{}, total int) {
	resp := &Response{
		Meta: &Meta{
			Total: uint(total),
		},
		Data: data,
	}
	Json(c, resp, http.StatusOK, false)
}

func Json(ctx *gin.Context, data interface{}, status int, abort bool) {
	if abort {
		ctx.AbortWithStatusJSON(status, data)
	} else {
		ctx.JSON(status, data)
	}
}
