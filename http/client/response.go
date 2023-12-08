package client

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

type Response struct {
	*http.Response
}

func (r *Response) IsError() bool {
	return r.Response.StatusCode >= 400
}

func (r *Response) ScanJsonBody(data any) (err error) {
	defer r.Response.Body.Close()
	content, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, data)
	return
}

func (r *Response) ScanXmlBody(data any) (err error) {
	defer r.Response.Body.Close()
	content, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return
	}
	err = xml.Unmarshal(content, data)
	return
}
