package util

import (
	"encoding/json"
	"github.com/kirinlabs/HttpRequest"
	"time"
)

// HTTPClient ...
type HTTPClient struct {
	headers map[string]string
	cookies map[string]string
}

// NewHTTPClient alloc tenant center handler
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

// SetHeader set heder
func (c *HTTPClient) SetHeader(key, val string) {
	c.headers[key] = val
}

// SetCookie set cookie
func (c *HTTPClient) SetCookie(key, val string) {
	c.cookies[key] = val
}

// Req get request handler
func (c *HTTPClient) Req() *HttpRequest.Request {
	req := HttpRequest.NewRequest().Debug(true).SetTimeout(3 * time.Second)
	req.SetHeaders(c.headers)
	req.SetCookies(c.cookies)
	return req
}

// Get ...
func (c *HTTPClient) Get(address string, v ...interface{}) (*HttpRequest.Response, error) {
	return c.Req().Get(address, v)
}

// Post ...
func (c *HTTPClient) Post(address string, v ...interface{}) (*HttpRequest.Response, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	req := c.Req()
	//fmt.Printf("post %s", string(data))
	res, err := req.JSON().Post(address, string(data))
	if err != nil {
		return nil, err
	}
	return res, nil
}
