package clientbak

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/zeddy-go/zeddy/mapx"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
)

type OptFunc func(*Client)

func WithBaseUrl(baseUrl string) OptFunc {
	return func(client *Client) {
		client.BaseUrl = baseUrl
	}
}

func WithTimeout(d time.Duration) OptFunc {
	return func(c *Client) {
		c.timeout = d
	}
}

func NewClient(opts ...OptFunc) *Client {
	c := &Client{
		header:  make(http.Header),
		query:   make(url.Values),
		cookies: make([]*http.Cookie, 0),
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

type Client struct {
	BaseUrl string
	debug   bool
	header  http.Header
	timeout time.Duration
	query   url.Values
	cookies []*http.Cookie
	clone   bool
}

func (c *Client) Clone() *Client {
	if c.clone {
		return c
	} else {
		return &Client{
			BaseUrl: c.BaseUrl,
			debug:   c.debug,
			timeout: c.timeout,
			header:  c.header.Clone(),
			query:   mapx.CloneSimpleMapSlice(c.query),
			cookies: cloneCookies(c.cookies),
			clone:   true,
		}
	}
}

func (c *Client) SetHeader(header http.Header) *Client {
	newClient := c.Clone()
	newClient.header = header
	return newClient
}

func (c *Client) AddHeader(key string, value string) *Client {
	newClient := c.Clone()
	newClient.header.Add(key, value)
	return newClient
}

func (c *Client) Debug() *Client {
	c.debug = true
	return c
}

func (c *Client) SetTimeout(d time.Duration) *Client {
	newClient := c.Clone()
	newClient.timeout = d
	return newClient
}

func (c *Client) SetQuery(values url.Values) *Client {
	newClient := c.Clone()
	newClient.query = values
	return newClient
}

func (c *Client) AddQuery(key string, value string) *Client {
	newClient := c.Clone()
	newClient.query.Add(key, value)
	return newClient
}

func (c *Client) SetCookies(cookies []*http.Cookie) *Client {
	newClient := c.Clone()
	newClient.cookies = cookies
	return newClient
}

func (c *Client) AddCookie(cookie *http.Cookie) *Client {
	newClient := c.Clone()
	newClient.cookies = append(c.cookies, cookie)
	return newClient
}

func (c *Client) Get(url string) (*Response, error) {
	return c.get(url)
}

func (c *Client) get(url string) (*Response, error) {
	return c.Request(c.makeRequest(http.MethodGet, url, nil))
}

func (c *Client) Delete(url string) (*Response, error) {
	return c.delete(url)
}

func (c *Client) delete(url string) (*Response, error) {
	return c.Request(c.makeRequest(http.MethodDelete, url, nil))
}

func (c *Client) PutJson(url string, data any) (resp *Response, err error) {
	content, err := json.Marshal(data)
	if err != nil {
		return
	}
	c.header.Set("Content-Type", "application/json")
	return c.put(url, bytes.NewReader(content))
}

func (c *Client) PutXml(url string, data any) (resp *Response, err error) {
	content, err := xml.Marshal(data)
	if err != nil {
		return
	}
	c.header.Set("Content-Type", "text/xml")
	return c.put(url, bytes.NewReader(content))
}

func (c *Client) put(url string, body io.Reader) (*Response, error) {
	return c.Request(c.makeRequest(http.MethodPut, url, body))
}

func (c *Client) PostJson(url string, data any) (resp *Response, err error) {
	content, err := json.Marshal(data)
	if err != nil {
		return
	}

	c.header.Set("Content-Type", "application/json")
	return c.post(url, bytes.NewReader(content))
}

func (c *Client) PostXml(url string, data any) (resp *Response, err error) {
	content, err := xml.Marshal(data)
	if err != nil {
		return
	}
	c.header.Set("Content-Type", "text/xml")
	return c.post(url, bytes.NewReader(content))
}

func (c *Client) PostForm(u string, data any) (resp *Response, err error) {
	c.header.Set("Content-Type", "application/x-www-form-urlencoded")
	vData := reflect.ValueOf(data)
	if vData.Kind() == reflect.Pointer {
		vData = vData.Elem()
	}
	if vData.Kind() == reflect.String {
		return c.post(u, strings.NewReader(vData.Interface().(string)))
	}

	if vData.Kind() != reflect.Struct {
		err = errors.New("unsupported type of data")
		return
	}

	tData := reflect.TypeOf(data)
	if tData.Kind() == reflect.Pointer {
		tData = tData.Elem()
	}

	values := url.Values{}
	for i := 0; i < vData.NumField(); i++ {
		fieldValue := vData.Field(i)
		if fieldValue.IsValid() {
			continue
		}
		fieldType := tData.Field(i)

		var (
			key   string
			value string
		)
		parts := strings.Split(fieldType.Tag.Get("form"), ",")
		if len(parts) > 0 {
			key = parts[0]
		} else {
			key = fieldType.Name
		}

		value = gconv.String(fieldValue.Interface())

		values.Add(key, value)
	}

	return c.post(u, strings.NewReader(values.Encode()))
}

func (c *Client) post(url string, body io.Reader) (*Response, error) {
	return c.Request(c.makeRequest(http.MethodPost, url, body))
}

func (c *Client) buildUrl(u string) string {
	if !strings.HasPrefix(u, "http") {
		u = fmt.Sprintf(
			"%s/%s",
			strings.TrimRight(c.BaseUrl, "/"),
			strings.TrimLeft(u, "/"),
		)
	}

	parsedUrl, err := url.Parse(u)
	if err != nil {
		panic(err)
	}

	for k, values := range c.query {
		for _, v := range values {
			q := parsedUrl.Query()
			q.Add(k, v)
			parsedUrl.RawQuery = q.Encode()
		}
	}

	return parsedUrl.String()
}

func (c *Client) makeRequest(method string, url string, body io.Reader) (req *http.Request) {
	req, err := http.NewRequest(method, c.buildUrl(url), body)
	if err != nil {
		panic(err)
	}

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	for k, values := range c.header {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	return
}

func (c *Client) Request(req *http.Request) (resp *Response, err error) {
	httpClient := http.Client{
		Timeout: c.timeout,
	}

	var (
		start   time.Time
		content []byte
	)

	if c.debug {
		if req.Body != nil {
			content, err = io.ReadAll(req.Body)
			if err != nil {
				return
			}
			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewReader(content))
		}
		start = time.Now()
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return
	}

	if c.debug {
		var (
			respContent []byte
		)
		respContent, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return
		}
		res.Body = io.NopCloser(bytes.NewReader(respContent))

		fmt.Printf(
			"[request: method=%s, url=%s, body=%s, header=%v]\n[request time: %fs]\n[response: body=%s]\n",
			req.Method,
			req.URL.String(),
			string(content),
			req.Header,
			time.Since(start).Seconds(),
			string(respContent),
		)
	}

	resp = &Response{
		Response: res,
	}

	return
}

func cloneCookies(src []*http.Cookie) (dest []*http.Cookie) {
	if src == nil {
		return nil
	}

	dest = make([]*http.Cookie, 0, len(src))
	for _, item := range src {
		tmp := &http.Cookie{
			Name:       item.Name,
			Value:      item.Value,
			Path:       item.Path,
			Domain:     item.Domain,
			Expires:    item.Expires,
			RawExpires: item.RawExpires,
			MaxAge:     item.MaxAge,
			Secure:     item.Secure,
			HttpOnly:   item.HttpOnly,
			SameSite:   item.SameSite,
			Raw:        item.Raw,
			Unparsed:   make([]string, len(item.Unparsed)),
		}
		copy(tmp.Unparsed, item.Unparsed)
		dest = append(dest, tmp)
	}

	return
}
