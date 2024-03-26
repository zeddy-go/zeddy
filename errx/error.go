package errx

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
)

func WithCode(code int) func(map[InfoKey]any) {
	return func(m map[InfoKey]any) {
		m[Code] = code
	}
}

func WithDetailMap(d map[string]any) func(map[InfoKey]any) {
	return func(m map[InfoKey]any) {
		m[Detail] = d
	}
}

func WithDetailSlice(d []any) func(map[InfoKey]any) {
	return func(m map[InfoKey]any) {
		m[Detail] = d
	}
}

func New(msg string, sets ...func(map[InfoKey]any)) error {
	return NewWithSkip(msg, 1, sets...)
}

func Wrap(err error, msg string, sets ...func(map[InfoKey]any)) error {
	return WrapWithSkip(err, msg, 1, sets...)
}

// NewFromStatus new an error from status.Status
func NewFromStatus(s *status.Status) error {
	e := NewWithSkip(s.Message(), 1)
	tmp := e.(*Errx)
	if statusDetails := s.Details(); len(statusDetails) > 0 {
		if details, ok := statusDetails[0].(*Map); ok {
			d := PbMap2MapStrAny(details)
			for key, value := range d {
				tmp.info[InfoKey(key)] = value
			}
		}
	}
	return tmp
}

func NewWithSkip(msg string, skip int, sets ...func(map[InfoKey]any)) error {
	e := make(map[InfoKey]any, 10)
	e[Msg] = msg

	_, file, line, _ := runtime.Caller(1 + skip)
	e[File] = file
	e[Line] = line

	e[ErrStack] = fmt.Sprintf(
		"%s:%d:%s\n",
		file,
		line,
		msg,
	)

	for _, set := range sets {
		set(e)
	}

	return &Errx{
		info: e,
	}
}

func WrapWithSkip(err error, message string, skip int, sets ...func(map[InfoKey]any)) error {
	if err == nil {
		return nil
	}
	e := NewWithSkip(message, skip+1, sets...).(*Errx)
	e.info[Err] = err

	switch x := err.(type) {
	case *Errx:
		e.info[ErrStack] = fmt.Sprintf(
			"%s:%d:%s\n%s",
			e.info[File],
			e.info[Line],
			e.info[Msg],
			x.ErrStack(),
		)
		if code, ok := x.info[Code]; ok && e.info[Code] == 0 {
			e.info[Code] = code
		}
	default:
		e.info[ErrStack] = fmt.Sprintf(
			"%s:%d:%s\n%s",
			e.info[File],
			e.info[Line],
			e.info[Msg],
			x.(error).Error(),
		)
	}

	return e
}

// Errx error扩展结构体
//
// note: 如果直接用map，则无法使用errors包的工具函数，如 errors.Is
type Errx struct {
	info map[InfoKey]any
}

func (e Errx) Error() string {
	msg := e.info[Msg].(string)
	if sub, ok := e.info[Err]; ok {
		switch x := sub.(type) {
		case Errx:
			msg += ": " + x.info[Msg].(string)
		default:
			msg += ": " + sub.(error).Error()
		}
	}

	return msg
}

func (e Errx) ErrStack() string {
	return e.info[ErrStack].(string)
}

func (e Errx) Unwrap() error {
	if err, ok := e.info[Err]; ok {
		return err.(error)
	} else {
		return nil
	}
}

type Fields map[InfoKey]any

func (e Errx) Set(fields Fields) {
	for key, value := range fields {
		e.info[key] = value
	}
}

func (e Errx) Get(field InfoKey) (value any, ok bool) {
	value, ok = e.info[field]
	return
}

func (e Errx) MustGet(field InfoKey) (value any) {
	value, _ = e.info[field]
	return
}

func (e Errx) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		switch {
		case s.Flag('+'):
			_, _ = s.Write([]byte(fmt.Sprintf("%s:%d:%s\n", e.info[File], e.info[Line], e.info[Msg])))
		case s.Flag('#'):
			_, _ = s.Write([]byte(e.info[ErrStack].(string)))
		default:
			if e.info[Err] != nil {
				_, _ = s.Write([]byte(fmt.Sprintf("%s]<=[%v\n", e.info[Msg], e.info[Err])))
			} else {
				_, _ = s.Write([]byte(e.info[Msg].(string) + "\n"))
			}
		}
	}
}

func (e Errx) GRPCStatus() *status.Status {
	var code codes.Code
	if e.info[Code] != nil {
		code = codes.Code(e.info[Code].(int))
	} else {
		code = codes.Unknown
	}
	s := status.New(code, e.Error())

	if len(e.info) > 0 {
		var (
			err error
		)
		m := &Map{}
		m.Fields = Map2Pb(e.info)
		s, err = s.WithDetails(m)
		if err != nil {
			panic(fmt.Errorf("build status failed: %w", err))
		}
	}

	return s
}

func (e Errx) HasCode() bool {
	_, ok := e.info[Code]
	return ok
}

func (e Errx) GetCode() (code int, ok bool) {
	c, ok := e.info[Code]
	if ok {
		code = c.(int)
	}

	return
}

func (e Errx) SetCode(code int) {
	e.info[Code] = code
}

// Is 判断 e 是否与 err 相等
//
// note: 与 errors.Is 的行为不同的是，**这个相等指各自包含的 code 相等**
func (e Errx) Is(err error) bool {
	for {
		if err == nil {
			return false
		}
		if x, ok := err.(*Errx); ok {
			code1, ok1 := e.info[Code]
			code2, ok2 := x.info[Code]
			msg1, ok3 := e.info[Msg]
			msg2, ok4 := x.info[Msg]
			if ((ok1 && ok2) && code1 == code2) || ((ok3 && ok4) && msg1 == msg2) {
				return true
			}
		}
		err = errors.Unwrap(err)
	}
}

func GetErrxField[T any](err any, key InfoKey) (result T) {
	tmp, _ := err.(*Errx).Get(key)
	if tmp != nil {
		return tmp.(T)
	}
	return
}

// Is 判断 err 是否与 target 相等
//
// note: 与 errors.Is 的行为不同的是，**这个相等还包括各自包含的 code 相等**
func Is(err error, target error) bool {
	if errors.Is(err, target) {
		return true
	}
	if x, ok := target.(*Errx); !ok {
		return false
	} else {
		return x.Is(err)
	}
}
