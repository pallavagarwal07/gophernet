// build: js

package main

import (
	"encoding/json"

	"github.com/gopherjs/gopherjs/js"
)

type Result struct {
	Name string
	Fail bool
	Outp string
}

func writeToDoc(results []Result) {
	out, err := json.Marshal(results)
	if err != nil {
		// Will never be hit as it depends only on type 'Result'
		// So we can risk using panic(). TODO: handle properly.
		panic(err)
	}
	window := js.Global.Get("window")
	window.Set("my_important_result", string(out))
}

func getVal(out string) string {
	var v interface{}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		return "--"
	}
	return v.(map[string]interface{})["value"].(string)
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
