package client

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

type iResponse interface {
	// Response return origin *http.Response, may close by yourself
	Response() (res *http.Response, err error)
	ScanJsonBody(data any) (err error)
	ScanXmlBody(data any) (err error)
}

func newResponse(origin *http.Response, err error) iResponse {
	return &response{
		origin: origin,
		err:    err,
	}
}

type response struct {
	origin *http.Response
	err    error
}

func (r response) Response() (res *http.Response, err error) {
	return r.origin, err
}

func (r response) ScanJsonBody(data any) (err error) {
	if r.err != nil {
		return r.err
	}
	defer r.origin.Body.Close()
	content, err := io.ReadAll(r.origin.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, data)
	return
}

func (r response) ScanXmlBody(data any) (err error) {
	if r.err != nil {
		return r.err
	}
	defer r.origin.Body.Close()
	content, err := io.ReadAll(r.origin.Body)
	if err != nil {
		return
	}
	err = xml.Unmarshal(content, data)
	return
}
