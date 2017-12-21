// +build js

package gophernet

import (
	"errors"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jsbuiltin"
	"honnef.co/go/js/xhr"
)

func GetIsNode() bool {
	typeOf := jsbuiltin.TypeOf

	if proc := js.Global.Get("process"); typeOf(proc) == "object" {
		if ver := proc.Get("versions"); typeOf(ver) == "object" {
			if node := ver.Get("node"); node != js.Undefined {
				return true
			}
		}
	}
	return false
}

func getErrorFunc(e chan error) func(map[string]interface{}, string) {
	return func(jqXHR map[string]interface{}, exception string) {
		var msg = ""
		status := int(jqXHR["status"].(float64))
		if status == 0 {
			msg = "Not connect.\n Verify Network."
		} else if status == 404 {
			msg = "Requested page not found. [404]"
		} else if status == 500 {
			msg = "Internal Server Error [500]."
		} else if exception == "parsererror" {
			msg = "Requested JSON parse failed."
		} else if exception == "timeout" {
			msg = "Time out error."
		} else if exception == "abort" {
			msg = "Ajax request aborted."
		} else {
			msg = "Uncaught Error.\n" + jqXHR["responseText"].(string)
		}
		e <- errors.New(msg)
	}
}

func (c *Client) get(urlStr string, params Values) ([]byte, error) {
	urlParsed, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	queriesInitial := urlParsed.Query()
	for key, arr := range params {
		for _, val := range arr {
			queriesInitial.Add(key, val)
		}
	}
	urlParsed.RawQuery = queriesInitial.Encode()
	urlStr = urlParsed.String()

	req := xhr.NewRequest("GET", urlStr)
	for key, list := range c.Header {
		for _, val := range list {
			req.SetRequestHeader(key, val)
		}
	}
	req.ResponseType = xhr.ArrayBuffer
	err = req.Send(nil)
	if err != nil {
		return nil, err
	}
	b := js.Global.Get("Uint8Array").New(req.Response).Interface().([]byte)
	return b, nil
}

func (c *Client) post(urlStr string, contentType string, body io.Reader) ([]byte, error) {
	req := xhr.NewRequest("POST", urlStr)
	req.ResponseType = xhr.ArrayBuffer
	for key, list := range c.Header {
		for _, val := range list {
			req.SetRequestHeader(key, val)
		}
	}
	req.SetRequestHeader("Content-Type", contentType)
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	if err := req.Send(string(data)); err != nil {
		return nil, err
	}
	b := js.Global.Get("Uint8Array").New(req.Response).Interface().([]byte)
	return b, nil
}
