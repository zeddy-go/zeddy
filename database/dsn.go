package database

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Type is database type. eg. mysql, postgres
type Type string

const (
	TypeMysql    Type = "mysql"
	TypePostgres Type = "postgres"
)

type DSN string

func (d DSN) Encode() string {
	return encodeQuery(string(d))
}

func (d DSN) split() []string {
	reg := regexp.MustCompile(`(^\S+)://(\S+$)`)
	info := reg.FindStringSubmatch(d.Encode())
	if len(info) < 3 {
		panic(errors.New("dsn is invalid, forget schema?"))
	}
	return info[1:]
}

func (d DSN) Type() Type {
	return Type(d.split()[0])
}

func (d DSN) RemoveSchema() string {
	return d.split()[1]
}

func encodeQuery(u string) string {
	tmp, err := url.Parse(u)
	if err == nil {
		return fmt.Sprintf("%s://%s@%s%s?%s", tmp.Scheme, tmp.User, tmp.Host, tmp.Path, tmp.Query().Encode())
	} else {
		urlArr := strings.Split(u, "?")
		if len(urlArr) > 1 {
			q, err := url.ParseQuery(urlArr[1])
			if err != nil {
				panic(err)
			}
			return fmt.Sprintf("%s?%s", urlArr[0], q.Encode())
		} else {
			return u
		}
	}
}
