// +build js

package netgopher

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

func init() {
	if k := js.Global.Get("jQuery"); k == js.Undefined {
		panic("FATAL: jQuery is not loaded.")
	}
}

func get(url string) string {
	c := make(chan string)
	jquery.Get(url, func(data string) {
		c <- data
	})
	return <-c
}
