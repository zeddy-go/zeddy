package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func NewResponse(r *http.Response) *Response {
	return &Response{
		Response: r,
	}
}

type Response struct {
	*http.Response
}

func (r *Response) ScanJson(v any) (err error) {
	content, err := r.GetBody()
	if err != nil {
		return
	}
	return json.Unmarshal(content, v)
}

func (r *Response) Info() (info map[string]any) {
	info = make(map[string]any)

	if r.Body != nil {
		body := r.Body
		defer body.Close()

		content, err := io.ReadAll(body)
		if err != nil {
			panic(err)
		}

		info["body"] = string(content)

		r.Body = io.NopCloser(bytes.NewReader(content))
	} else {
		info["body"] = "null"
	}

	return
}

func (r *Response) GetBody() (body []byte, err error) {
	defer r.Body.Close()
	body, err = io.ReadAll(r.Body)
	if err != nil {
		return
	}
	return
}
