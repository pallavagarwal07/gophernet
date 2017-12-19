package gophernet

import (
	"io"
	"net/url"
)

var Get func(string, url.Values) ([]byte, error) = get
var PostForm func(string, url.Values) ([]byte, error) = postform
var Post func(url string, contentType string, body io.Reader) ([]byte, error) = post
