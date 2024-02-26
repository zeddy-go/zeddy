package client

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func WithBaseUrl(baseUrl string) func(*Client) {
	return func(client *Client) {
		client.baseUrl = baseUrl
	}
}

func NewClient(opts ...func(*Client)) *Client {
	c := &Client{
		cookies: make([]*http.Cookie, 0),
		headers: make(http.Header),
		queries: make(url.Values),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type Client struct {
	cookies []*http.Cookie
	headers http.Header
	clone   int
	queries url.Values
	debug   bool
	baseUrl string
}

func (c *Client) BaseUrl(baseUrl string) (client *Client) {
	client = c.Clone()
	client.baseUrl = strings.TrimRight(baseUrl, "/")
	return
}

func (c *Client) Debug() (client *Client) {
	client = c.Clone()
	client.debug = true
	return
}

func (c *Client) AddCookies(cookies ...*http.Cookie) (client *Client) {
	client = c.Clone()
	client.cookies = append(client.cookies, cookies...)
	return
}

func (c *Client) SetHeaders(header map[string][]string) (client *Client) {
	client = c.Clone()
	client.headers = header
	return
}

func (c *Client) AddHeader(key string, values ...string) (client *Client) {
	client = c.Clone()
	for _, value := range values {
		client.headers.Add(key, value)
	}

	return
}

func (c *Client) SetHeader(key string, values ...string) (client *Client) {
	client = c.Clone()
	for idx, value := range values {
		if idx == 0 {
			client.headers.Set(key, value)
		} else {
			client.headers.Add(key, value)
		}
	}

	return
}

func (c *Client) SetQueries(queries map[string][]string) (client *Client) {
	client = c.Clone()
	client.queries = queries
	return
}

func (c *Client) AddQuery(key string, values ...string) (client *Client) {
	client = c.Clone()
	for _, value := range values {
		client.queries.Add(key, value)
	}

	return
}

func (c *Client) SetQuery(key string, values ...string) (client *Client) {
	client = c.Clone()
	for idx, value := range values {
		if idx == 0 {
			client.queries.Set(key, value)
		} else {
			client.queries.Add(key, value)
		}
	}

	return
}

func (c *Client) Clone() *Client {
	switch c.clone {
	case 0:
		return &Client{
			cookies: c.cookies,
			headers: c.headers,
			clone:   1,
			queries: c.queries,
			debug:   c.debug,
			baseUrl: c.baseUrl,
		}
	case 1:
		return c
	default:
		return c
	}
}

func (c *Client) Clean() {
	c.cookies = make([]*http.Cookie, 0)
	c.headers = make(http.Header)
}

func (c *Client) makeRequest(method string, url string, body any) (req *Request, err error) {
	return NewRequest(method, url, body, c.cookies, c.headers)
}

func (c *Client) Do(request *Request) (resp *Response, err error) {
	response, err := http.DefaultClient.Do(request.Request)
	if err != nil {
		return
	}
	resp = NewResponse(response)
	return
}

func (c *Client) buildUrl(uri string) (u string, err error) {
	var tmp *url.URL
	if c.baseUrl != "" {
		tmp, err = url.Parse(fmt.Sprintf("%s/%s", c.baseUrl, strings.TrimLeft(uri, "/")))
	} else {
		tmp, err = url.Parse(uri)
	}
	if err != nil {
		return
	}

	if len(c.queries) > 0 {
		tmp.RawQuery = c.queries.Encode()
	}

	u = tmp.String()
	return
}

func (c *Client) PostJson(uri string, body any) (resp *Response, err error) {
	c.SetHeader("Content-Type", "application/json")
	return c.Post(uri, body)
}

func (c *Client) Post(uri string, body any) (resp *Response, err error) {
	u, err := c.buildUrl(uri)
	if err != nil {
		return
	}

	req, err := c.makeRequest(http.MethodPost, u, body)
	if err != nil {
		return
	}

	resp, err = c.Do(req)

	if c.debug {
		slog.Debug("request debug", "request", req.Info(), "response", resp.Info())
	}

	return
}

func (c *Client) PutJson(uri string, body any) (resp *Response, err error) {
	c.SetHeader("Content-Type", "application/json")
	return c.Put(uri, body)
}

func (c *Client) Put(uri string, body any) (resp *Response, err error) {
	u, err := c.buildUrl(uri)
	if err != nil {
		return
	}

	req, err := c.makeRequest(http.MethodPut, u, body)
	if err != nil {
		return
	}

	resp, err = c.Do(req)

	if c.debug {
		slog.Debug("request debug", "request", req.Info(), "response", resp.Info())
	}

	return
}

func (c *Client) Delete(uri string) (resp *Response, err error) {
	u, err := c.buildUrl(uri)
	if err != nil {
		return
	}

	req, err := c.makeRequest(http.MethodDelete, u, nil)
	if err != nil {
		return
	}

	resp, err = c.Do(req)

	if c.debug {
		slog.Debug("request debug", "request", req.Info(), "response", resp.Info())
	}

	return
}

func (c *Client) Get(uri string) (resp *Response, err error) {
	u, err := c.buildUrl(uri)
	if err != nil {
		return
	}

	req, err := c.makeRequest(http.MethodGet, u, nil)
	if err != nil {
		return
	}

	resp, err = c.Do(req)

	if c.debug {
		slog.Debug("request debug", "request", req.Info(), "response", resp.Info())
	}

	return
}
