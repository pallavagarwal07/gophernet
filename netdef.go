// +build !js

package gophernet

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

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

	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func post(urlStr string, params url.Values) ([]byte, error) {
	resp, err := http.PostForm(urlStr, params)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
