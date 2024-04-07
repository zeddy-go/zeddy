package ginx

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zeddy-go/zeddy/errx"
	"gorm.io/gorm"
	"net/http"
)

func WithFileName(name string) func(*File) {
	return func(file *File) {
		file.name = name
	}
}

func WithMimeType(mimeType string) func(*File) {
	return func(file *File) {
		file.mimeType = mimeType
	}
}

func NewFile(content []byte, opts ...func(*File)) *File {
	f := &File{
		content: content,
	}

	for _, opt := range opts {
		opt(f)
	}

	if f.mimeType == "" {
		f.mimeType = http.DetectContentType(content)
	}

	return f
}

type File struct {
	name     string
	mimeType string
	content  []byte
}

func (f *File) Name() string {
	return f.name
}

func (f *File) MimeType() string {
	return f.mimeType
}

func (f *File) Content() []byte {
	return f.content
}

type IFile interface {
	Name() string
	MimeType() string
	Content() []byte
}

type IMeta interface {
	GetMeta() map[string]any
}

type Meta struct {
	Code        int  `json:"code,omitempty"`
	CurrentPage uint `json:"currentPage,omitempty"`
	Total       uint `json:"total,omitempty"`
	LastPage    uint `json:"lastPage,omitempty"`
	PerPage     uint `json:"perPage,omitempty"`
}

func (m *Meta) GetMeta() (result map[string]any) {
	s, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(s, &result)
	if err != nil {
		panic(err)
	}
	return
}

type Response struct {
	Data    interface{}    `json:"data"`
	Message string         `json:"message"`
	Meta    map[string]any `json:"meta,omitempty"`
}

type ResponseWithCode struct {
	Code int `json:"code"`
	Response
}

func NewAllOkStatusJsonResponse(data any, meta IMeta, err error) IResponse[*gin.Context] {
	resp := &AllOkStatusJsonResponse{
		Err: err,
		Response: &ResponseWithCode{
			Response: Response{},
		},
		DefaultErrCode: 1,
	}
	if data != nil {
		resp.Response.Data = data
	}
	if meta != nil {
		resp.Response.Meta = meta.GetMeta()
		if code, ok := resp.Response.Meta["code"]; ok {
			resp.Response.Code = int(code.(float64))
			delete(resp.Response.Meta, "code")
		}
	}

	return resp
}

type AllOkStatusJsonResponse struct {
	Err            error
	Response       *ResponseWithCode
	DefaultErrCode int
}

func (r *AllOkStatusJsonResponse) Do(ctx *gin.Context) {
	if r.Err != nil {
		var abort bool
		var data any
		message := r.Err.Error()
		status := http.StatusOK
		r.Response.Code = r.DefaultErrCode
		if x, ok := r.Err.(*errx.Errx); ok {
			if c, ok := x.Get(errx.Code); ok {
				r.Response.Code = c.(int)
			}
			abort = errx.GetErrxField[bool](x, errx.Abort)
		} else if _, ok := r.Err.(validator.ValidationErrors); ok {
			//code = http.StatusUnprocessableEntity
			//TODO: i18n and detail
		} else if errors.Is(r.Err, gorm.ErrRecordNotFound) {
			//code = http.StatusNotFound
		}
		r.Response.Response = Response{
			Message: message,
			Data:    data,
		}
		if abort {
			ctx.AbortWithStatusJSON(status, r.Response)
		} else {
			ctx.JSON(status, r.Response)
		}
		return
	}

	if r.Response == nil || r.Response.Data == nil {
		ctx.JSON(http.StatusOK, nil)
		return
	}

	ctx.JSON(http.StatusOK, r.Response)
}

func NewJsonResponse(data any, meta IMeta, err error) IResponse[*gin.Context] {
	resp := &JsonResponse{
		Err: err,
	}
	if data != nil || meta != nil {
		rr := &Response{}
		if data != nil {
			rr.Data = data
		}
		if meta != nil {
			rr.Meta = meta.GetMeta()
		}
		resp.Response = rr
	}

	return resp
}

type JsonResponse struct {
	Err      error
	Response *Response
}

func (r *JsonResponse) Do(ctx *gin.Context) {
	if r.Err != nil {
		var abort bool
		var data any
		message := r.Err.Error()
		status := http.StatusInternalServerError
		if x, ok := r.Err.(*errx.Errx); ok && http.StatusText(errx.GetErrxField[int](x, errx.Code)) != "" {
			status = errx.GetErrxField[int](x, errx.Code)
			abort = errx.GetErrxField[bool](x, errx.Abort)
		} else if _, ok := r.Err.(validator.ValidationErrors); ok {
			status = http.StatusUnprocessableEntity
			//TODO: i18n and detail
		} else if errors.Is(r.Err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		r.Response = &Response{
			Message: message,
			Data:    data,
		}
		if abort {
			ctx.AbortWithStatusJSON(status, r.Response)
		} else {
			ctx.JSON(status, r.Response)
		}
		return
	}

	if r.Response == nil || r.Response.Data == nil {
		ctx.JSON(http.StatusNoContent, nil)
		return
	}

	ctx.JSON(http.StatusOK, r.Response)
}

type IResponse[CTX any] interface {
	Do(CTX)
}
