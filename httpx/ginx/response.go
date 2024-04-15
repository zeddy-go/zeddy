package ginx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zeddy-go/zeddy/convert"
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
	Code        int
	CurrentPage uint
	Total       uint
	LastPage    uint
	PerPage     uint
}

func (m *Meta) GetMeta() (result map[string]any) {
	result = make(map[string]any)
	if m.Code != 0 {
		result["code"] = m.Code
	}
	if m.CurrentPage != 0 {
		result["currentPage"] = m.CurrentPage
	}
	if m.Total != 0 {
		result["total"] = m.Total
	}
	if m.LastPage != 0 {
		result["lastPage"] = m.LastPage
	}
	if m.PerPage != 0 {
		result["perPage"] = m.PerPage
	}
	return
}

func NewRestfulResponse() IResponse[*gin.Context] {
	return &RestfulResponse{}
}

type RestfulResponse struct {
	Data    any            `json:"data"`
	Message string         `json:"message"`
	Meta    map[string]any `json:"meta,omitempty"`
	file    IFile
	meta    IMeta
	err     error
}

func (r *RestfulResponse) SetData(data any) IResponse[*gin.Context] {
	r.Data = data
	return r
}

func (r *RestfulResponse) SetMeta(meta IMeta) IResponse[*gin.Context] {
	r.meta = meta
	return r
}

func (r *RestfulResponse) SetError(err error) IResponse[*gin.Context] {
	r.err = err
	return r
}

func (r *RestfulResponse) SetFile(file IFile) IResponse[*gin.Context] {
	r.file = file
	return r
}

func (r *RestfulResponse) Do(ctx *gin.Context) {
	if r.err != nil {
		r.Data = nil
		var abort bool
		message := r.err.Error()
		status := http.StatusInternalServerError
		if x, ok := r.err.(*errx.Errx); ok && http.StatusText(errx.GetErrxField[int](x, errx.Code)) != "" {
			status = errx.GetErrxField[int](x, errx.Code)
			abort = errx.GetErrxField[bool](x, errx.Abort)
		} else if _, ok := r.err.(validator.ValidationErrors); ok {
			status = http.StatusUnprocessableEntity
			//TODO: i18n and detail
		} else if errors.Is(r.err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		r.Message = message
		if abort {
			ctx.AbortWithStatusJSON(status, r)
		} else {
			ctx.JSON(status, r)
		}
		return
	}

	if r.file != nil {
		ctx.Status(http.StatusOK)
		if r.file.MimeType() != "" {
			ctx.Writer.Header().Add("Content-Type", r.file.MimeType())
		} else {
			ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
		}
		if r.file.Name() != "" {
			ctx.Writer.Header().Add("Content-Disposition", "attachment; filename=\""+r.file.Name()+"\"")
		}
		_, _ = ctx.Writer.Write(r.file.Content())
		return
	}

	var status int
	if r.Data == nil {
		status = http.StatusNoContent
	} else {
		status = http.StatusOK
	}

	if r.meta != nil {
		r.Meta = r.meta.GetMeta()
	}
	if code, ok := r.Meta["code"]; ok {
		delete(r.Meta, "code")
		if c, e := convert.To[int](code); e == nil && http.StatusText(c) != "" {
			status = c
		}
	}
	ctx.JSON(status, r)
}

type IResponse[CTX any] interface {
	SetData(any) IResponse[CTX]
	SetMeta(IMeta) IResponse[CTX]
	SetError(error) IResponse[CTX]
	SetFile(IFile) IResponse[CTX]
	Do(CTX)
}
