package gophernet

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Values = url.Values
type Header = http.Header

type Client struct {
	Header Header
}

var defClient = Client{Header: nil}

var Get func(string, Values) ([]byte, error) = defClient.Get
var PostForm func(string, Values) ([]byte, error) = defClient.PostForm
var Post func(url string, contentType string, body io.Reader) ([]byte, error) = defClient.Post

func (c *Client) Get(s string, v Values) ([]byte, error) {
	return c.get(s, v)
}
func (c *Client) PostForm(urlStr string, data Values) ([]byte, error) {
	return c.Post(urlStr, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
func (c *Client) Post(s string, u string, b io.Reader) ([]byte, error) {
	return c.post(s, u, b)
}
