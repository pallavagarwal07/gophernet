// +build !js

package gophernet

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

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

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	for key, list := range c.Header {
		for _, val := range list {
			req.Header[key] = append(req.Header[key], val)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) post(url string, contentType string, body io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	for key, list := range c.Header {
		for _, val := range list {
			req.Header[key] = append(req.Header[key], val)
		}
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return output, nil
}
