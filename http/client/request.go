package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func NewRequest(method string, url string, body any, cookies []*http.Cookie, headers http.Header) (req *Request, err error) {
	req = &Request{}

	var b io.Reader
	if body != nil {
		switch x := body.(type) {
		case string:
			b = strings.NewReader(x)
			req.bodyCache = x
		case []byte:
			b = bytes.NewReader(x)
			req.bodyCache = string(x)
		default:
			var byts []byte
			byts, err = json.Marshal(body)
			if err != nil {
				return
			}
			b = bytes.NewBuffer(byts)
			req.bodyCache = string(byts)
		}
	}

	r, err := http.NewRequest(method, url, b)
	if err != nil {
		return
	}

	if len(cookies) > 0 {
		for _, cookie := range cookies {
			r.AddCookie(cookie)
		}
	}

	if len(headers) > 0 {
		r.Header = headers
	}

	req.Request = r

	return req, nil
}

type Request struct {
	*http.Request
	bodyCache string
}

func (r *Request) Info() map[string]any {
	return map[string]any{
		"method": r.Method,
		"url":    r.URL.String(),
		"header": r.Header,
		"body":   r.bodyCache,
	}
}
