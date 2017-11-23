// build: js

package main

import (
	"encoding/json"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)

type customT struct {
	name    string
	failed  bool
	skipped bool
	skipStr string
	log     string
	err     string
}

var _ TB = &customT{}

func (*customT) Helper() {
	return
}

func (c *customT) Skipped() bool {
	return c.skipped
}

func (c *customT) Skipf(format string, args ...interface{}) {
	c.skipped = true
	c.skipStr = fmt.Sprintf(format, args...)
}

func (c *customT) SkipNow() {
	c.skipped = true
}

func (c *customT) Skip(args ...interface{}) {
	c.skipped = true
	c.skipStr = fmt.Sprint(args...)
}

func (c *customT) Name() string {
	return c.name
}

func (c *customT) Log(args ...interface{}) {
	c.log += "\n" + fmt.Sprint(args...)
}

func (c *customT) Logf(format string, args ...interface{}) {
	c.log += "\n" + fmt.Sprintf(format, args...)
}

func (c *customT) Fail() {
	c.failed = true
}

func (c *customT) FailNow() {
	c.failed = true
}

func (c *customT) Failed() bool {
	return c.failed
}

func (c *customT) Error(args ...interface{}) {
	c.failed = true
	c.err = fmt.Sprint(args...)
}

func (c *customT) Errorf(format string, args ...interface{}) {
	c.failed = true
	c.err = fmt.Sprintf(format, args...)
}

func (c *customT) Fatal(args ...interface{}) {
	c.failed = true
	c.err = fmt.Sprint(args...)
}

func (c *customT) Fatalf(format string, args ...interface{}) {
	c.failed = true
	c.err = fmt.Sprintf(format, args...)
}

type Result struct {
	Name string
	Fail bool
	Outp string
}

func main() {
	results := []Result{}
	for name, test := range tests {
		c := &customT{
			name:    name,
			failed:  false,
			skipped: false,
			skipStr: "",
			log:     "",
			err:     "",
		}
		test(c)
		r := &Result{
			Name: name,
			Fail: c.Failed(),
			Outp: c.err,
		}
		results = append(results, *r)
	}
	writeToDoc(results)
}

func writeToDoc(results []Result) {
	out, err := json.Marshal(results)
	if err != nil {
		// Will never be hit as it depends only on type 'Result'
		// So we can risk using panic(). TODO: handle properly.
		panic(err)
	}
	document := js.Global.Get("document")
	div := document.Call("createElement", "div")
	div.Call("setAttribute", "id", "output_test")
	div.Set("innerHTML", string(out))
	document.Get("head").Call("appendChild", div)
}
