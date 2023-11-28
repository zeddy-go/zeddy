package errx

import "github.com/zeddy-go/zeddy/slicex"

type DetailKey string

const (
	Msg      DetailKey = "msg"
	Err      DetailKey = "err"
	File     DetailKey = "file"
	Line     DetailKey = "line"
	ErrStack DetailKey = "errStack"
	Code     DetailKey = "code"
	Detail   DetailKey = "detail"
)

var reserveKeys = []DetailKey{
	Msg,
	Err,
	File,
	Line,
	ErrStack,
	Code,
	Detail,
}

func (d DetailKey) IsReserved() bool {
	return slicex.Contains(d, reserveKeys)
}
