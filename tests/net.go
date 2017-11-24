package main

import (
	"encoding/json"

	"github.com/pallavagarwal07/gophernet"
)

// We can't use the standard Go test toolchain because we need to run the same
// tests with phantomjs + gopherjs compiled js files. Renaming this file with
// *_test.go format causes the `go run` to refuse to compile the test_runner.go
// (used for the js tests).
var tests = map[string]func(t TB){
	"TestGet":        TestGet,
	"TestGet_Method": TestGet_Method,
	"TestPost":       TestPost,
}

func TestGet(t TB) {
	got, err := gophernet.Get("http://localhost:"+PORT, nil)
	if err != nil {
		t.Fatalf("GET failed with error %v", err)
	}
	if want := "Hello World!"; string(got) != want {
		t.Fatalf("Got output: %q, Want: %q", string(got), want)
	}
}

func TestGet_Method(t TB) {
	got, err := gophernet.Get("http://localhost:"+PORT+"/echo", nil)
	if err != nil {
		t.Fatalf("GET failed with error %v", err)
	}
	var output Request
	if err := json.Unmarshal(got, &output); err != nil {
		t.Fatalf("JSON Unmarshal failed with error %v", err)
	}
	if want := "GET"; output.Method != want {
		t.Errorf("Incorrect method: Got %q, Want %q", output.Method, want)
	}
	if want := 0; len(output.Params) != want {
		t.Errorf("Incorrect params: Got length %d, Want %d", len(output.Params), want)
	}
}

func TestPost(t TB) {
	got, err := gophernet.Post("http://localhost:"+PORT, nil)
	if err != nil {
		t.Fatalf("POST failed with error %v", err)
	}
	if want := "Hello World!"; string(got) != want {
		t.Fatalf("Got output: %q, Want: %q", string(got), want)
	}
}
