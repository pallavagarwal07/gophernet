package main

import (
	"github.com/pallavagarwal07/gophernet"
)

// We can't use the standard Go test toolchain because we need to run the same
// tests with phantomjs + gopherjs compiled js files. Renaming this file with
// *_test.go format causes the `go run` to refuse to compile the test_runner.go
// (used for the js tests).
var tests = map[string]func(t TB){
	"TestGet":  TestGet,
	"TestPost": TestPost,
}

func TestGet(t TB) {
	got, err := gophernet.Get("http://localhost:"+PORT, nil)
	if err != nil {
		t.Fatalf("Get failed with error %v", err)
	}
	if want := "Hello World!"; string(got) != want {
		t.Fatalf("Got output: %q, Want: %q", string(got), want)
	}
}

func TestPost(t TB) {
	// Pass for now
}
