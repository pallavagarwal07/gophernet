// build: js

package main

import (
	"encoding/json"
	"fmt"

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

func runOneTest(name string, test func(t TB)) (result *Result) {
	result = &Result{}
	c := &customT{
		name:    name,
		failed:  false,
		skipped: false,
		skipStr: "",
		log:     "",
		err:     "",
	}

	defer func() {
		result.Name = name
		result.Fail = c.Failed()
		result.Outp = c.err
		if r := recover(); r != nil {
			result.Fail = true
		}
	}()

	t := &myT{
		TB:     c,
		backup: nil,
	}
	c.parent = t
	test(t)

	return result
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in main", r)
		}
	}()
	results := []Result{}
	for name, test := range tests {
		result := runOneTest(name, test)
		results = append(results, *result)
	}
	writeToDoc(results)
}
