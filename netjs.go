// +build js

package gophernet

import (
	"errors"
	"net/url"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/gopherjs/jsbuiltin"
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

func init() {
	// Check if jquery is loaded, and load if not. Note that the automatic
	// loading is needed for compatibility with `gopherjs serve`.
	if k := js.Global.Get("jQuery"); k == js.Undefined {
		c := make(chan string)
		callback := func() {
			c <- "done"
		}
		js.Global.Set("callbackFn", callback)

		srcCDN := "https://code.jquery.com/jquery-3.2.1.min.js"
		document := js.Global.Get("document")
		s := document.Call("createElement", "script")
		s.Call("setAttribute", "src", srcCDN)
		s.Call("setAttribute", "type", "text/javascript")
		s.Call("setAttribute", "onload", "callbackFn()")
		document.Get("head").Call("appendChild", s)

		// Stall the init function till jquery has been loaded.
		<-c
	}
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

func get(urlStr string, params url.Values) ([]byte, error) {
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

	d := make(chan string)
	e := make(chan error)
	jquery.Ajax(map[string]interface{}{
		"type": "GET",
		"url":  urlStr,
		"data": []string{},
		"success": func(data string) {
			d <- data
		},
		"error": getErrorFunc(e),
	})
	select {
	case out := <-d:
		return []byte(out), nil
	case out := <-e:
		return nil, out
	}
}

func post(urlStr string, params url.Values) ([]byte, error) {
	d := make(chan string)
	e := make(chan error)
	jquery.Ajax(map[string]interface{}{
		"type":        "POST",
		"url":         urlStr,
		"data":        params,
		"traditional": "true",
		"success": func(data string) {
			d <- data
		},
		"error": getErrorFunc(e),
	})
	select {
	case out := <-d:
		return []byte(out), nil
	case out := <-e:
		return nil, out
	}
}
