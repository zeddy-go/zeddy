package errx

import "github.com/zeddy-go/core/slicex"

type DetailKey string

const (
	Msg      DetailKey = "msg"
	Err      DetailKey = "err"
	File     DetailKey = "file"
	Line     DetailKey = "line"
	ErrStack DetailKey = "errStack"
	Code     DetailKey = "code"
)

var reserveKeys = []DetailKey{
	Msg,
	Err,
	File,
	Line,
	ErrStack,
	Code,
}

func (d DetailKey) IsReserved() bool {
	return slicex.Contains(d, reserveKeys)
}
