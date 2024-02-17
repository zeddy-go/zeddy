package errx

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
)

func WithCode(code int) func(map[DetailKey]any) {
	return func(m map[DetailKey]any) {
		m[Code] = code
	}
}

func WithDetailMap(d map[string]any) func(map[DetailKey]any) {
	return func(m map[DetailKey]any) {
		m[Detail] = d
	}
}

func WithDetailSlice(d []any) func(map[DetailKey]any) {
	return func(m map[DetailKey]any) {
		m[Detail] = d
	}
}

func New(msg string, sets ...func(map[DetailKey]any)) error {
	return NewWithSkip(msg, 1, sets...)
}

func Wrap(err error, msg string, sets ...func(map[DetailKey]any)) error {
	if err == nil {
		return nil
	}
	return WrapWithSkip(err, msg, 1, sets...)
}

// NewFromStatus new an error from status.Status
func NewFromStatus(s *status.Status) error {
	e := NewWithSkip(s.Message(), 1)
	tmp := e.(Errx)
	if statusDetails := s.Details(); len(statusDetails) > 0 {
		if details, ok := statusDetails[0].(*Map); ok {
			d := PbMap2MapStrAny(details)
			for key, value := range d {
				tmp[DetailKey(key)] = value
			}
		}
	}
	return tmp
}

func NewWithSkip(msg string, skip int, sets ...func(map[DetailKey]any)) error {
	e := make(map[DetailKey]any, 10)
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

	return Errx(e)
}

func WrapWithSkip(err error, message string, skip int, sets ...func(map[DetailKey]any)) error {
	e := NewWithSkip(message, skip+1, sets...)
	tmp := e.(Errx)
	tmp[Err] = err

	switch x := err.(type) {
	case *Errx:
		tmp[ErrStack] = fmt.Sprintf(
			"%s:%d:%s\n%s",
			tmp[File],
			tmp[Line],
			tmp[Msg],
			x.ErrStack(),
		)
	case Errx:
		tmp[ErrStack] = fmt.Sprintf(
			"%s:%d:%s\n%s",
			tmp[File],
			tmp[Line],
			tmp[Msg],
			x.ErrStack(),
		)
	default:
		tmp[ErrStack] = fmt.Sprintf(
			"%s:%d:%s\n%s",
			tmp[File],
			tmp[Line],
			tmp[Msg],
			x.(error).Error(),
		)
	}

	return tmp
}

type Errx map[DetailKey]any

func (e Errx) Error() string {
	msg := e[Msg].(string)
	if sub, ok := e[Err]; ok {
		switch x := sub.(type) {
		case *Errx:
			msg += ":" + (*x)[Msg].(string)
		case Errx:
			msg += ":" + x[Msg].(string)
		default:
			msg += ":" + sub.(error).Error()
		}
	}

	return msg
}

func (e Errx) ErrStack() string {
	return e[ErrStack].(string)
}

type Fields map[DetailKey]any

func (e Errx) Set(fields Fields) {
	for key, value := range fields {
		e[key] = value
	}
}

func (e Errx) Get(field DetailKey) (value any, ok bool) {
	value, ok = e[field]
	return
}

func (e Errx) Format(s fmt.State, c rune) {
	switch c {
	case 'v':
		switch {
		case s.Flag('+'):
			_, _ = s.Write([]byte(fmt.Sprintf("%s:%d:%s\n", e[File], e[Line], e[Msg])))
		case s.Flag('#'):
			_, _ = s.Write([]byte(e[ErrStack].(string)))
		default:
			if e[Err] != nil {
				_, _ = s.Write([]byte(fmt.Sprintf("%s]<=[%v\n", e[Msg], e[Err])))
			} else {
				_, _ = s.Write([]byte(e[Msg].(string) + "\n"))
			}
		}
	}
}

func (e Errx) GRPCStatus() *status.Status {
	var code codes.Code
	if e[Code] != nil {
		code = codes.Code(e[Code].(int))
	} else {
		code = codes.Unknown
	}
	s := status.New(code, e.Error())

	if len(e) > 0 {
		var (
			err error
		)
		m := &Map{}
		m.Fields = Map2Pb(e)
		s, err = s.WithDetails(m)
		if err != nil {
			panic(fmt.Errorf("build status failed: %w", err))
		}
	}

	return s
}

func GetErrxField[T any](err any, key DetailKey) (result T) {
	tmp, _ := err.(Errx).Get(key)
	if tmp != nil {
		return tmp.(T)
	}
	return
}
