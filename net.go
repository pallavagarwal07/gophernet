package gophernet

import "net/url"

var Get func(string, url.Values) ([]byte, error) = get
var Post func(string, url.Values) ([]byte, error) = post
