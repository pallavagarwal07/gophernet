package main

import (
	"encoding/json"
	"net/url"
	"reflect"

	"github.com/pallavagarwal07/gophernet"
)

// We can't use the standard Go test toolchain because we need to run the same
// tests with phantomjs + gopherjs compiled js files. Renaming this file with
// *_test.go format causes the `go run` to refuse to compile the test_runner.go
// (used for the js tests).
var tests = map[string]func(t TB){
	"TestGet":            TestGet,
	"TestGet_Method":     TestGet_Method,
	"TestGet_Binary":     TestGet_Binary,
	"TestGet_ParamsURL":  TestGet_ParamsURL,
	"TestGet_ParamsArg":  TestGet_ParamsArg,
	"TestGet_ParamsMix":  TestGet_ParamsMix,
	"TestPost":           TestPost,
	"TestPost_Method":    TestPost_Method,
	"TestPost_Binary":    TestPost_Binary,
	"TestPost_ParamsURL": TestPost_ParamsURL,
	"TestPost_ParamsArg": TestPost_ParamsArg,
	"TestPost_ParamsMix": TestPost_ParamsMix,
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

func TestGet_Binary(t TB) {
	got, err := gophernet.Get("http://localhost:"+PORT+"/binary", nil)
	if err != nil {
		t.Fatalf("GET failed with error %v", err)
	}
	want := deterministicBinData()
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Binary blob data does not match: Got %v, Want %v", got, want)
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

func TestGet_ParamsURL(t TB) {
	got, err := gophernet.Get("http://localhost:"+PORT+"/echo?a=b&c=d&a=g", nil)
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
	wantParams := []Pair{{"a", "b"}, {"a", "g"}, {"c", "d"}}
	if !reflect.DeepEqual(wantParams, output.Params) {
		t.Errorf("Incorrect params: Got %v, Want %v", output.Params, wantParams)
	}
}

func TestGet_ParamsArg(t TB) {
	args := url.Values{
		"c": {"d"},
		"a": {"b", "g"},
	}
	got, err := gophernet.Get("http://localhost:"+PORT+"/echo", args)
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
	wantParams := []Pair{{"a", "b"}, {"a", "g"}, {"c", "d"}}
	if !reflect.DeepEqual(wantParams, output.Params) {
		t.Errorf("Incorrect params: Got %v, Want %v", output.Params, wantParams)
	}
}

func TestGet_ParamsMix(t TB) {
	args := url.Values{
		"c": {"d"},
		"a": {"b", "g"},
	}
	got, err := gophernet.Get("http://localhost:"+PORT+"/echo?a=l&m=n", args)
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
	wantParams := []Pair{{"a", "b"}, {"a", "g"}, {"a", "l"}, {"c", "d"}, {"m", "n"}}
	if !reflect.DeepEqual(wantParams, output.Params) {
		t.Errorf("Incorrect params: Got %v, Want %v", output.Params, wantParams)
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

func TestPost_Binary(t TB) {
	got, err := gophernet.Post("http://localhost:"+PORT+"/binary", nil)
	if err != nil {
		t.Fatalf("POST failed with error %v", err)
	}
	want := deterministicBinData()
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Binary blob data does not match: Got %v, Want %v", got, want)
	}
}

func TestPost_Method(t TB) {
	got, err := gophernet.Post("http://localhost:"+PORT+"/echo", nil)
	if err != nil {
		t.Fatalf("POST failed with error %v", err)
	}
	var output Request
	if err := json.Unmarshal(got, &output); err != nil {
		t.Fatalf("JSON Unmarshal failed with error %v", err)
	}
	if want := "POST"; output.Method != want {
		t.Errorf("Incorrect method: Got %q, Want %q", output.Method, want)
	}
	if want := 0; len(output.Params) != want {
		t.Errorf("Incorrect params: Got length %d, Want %d", len(output.Params), want)
	}
}

func TestPost_ParamsURL(t TB) {
	got, err := gophernet.Post("http://localhost:"+PORT+"/echo?a=b&c=d&a=g", nil)
	if err != nil {
		t.Fatalf("POST failed with error %v", err)
	}
	var output Request
	if err := json.Unmarshal(got, &output); err != nil {
		t.Fatalf("JSON Unmarshal failed with error %v", err)
	}
	if want := "POST"; output.Method != want {
		t.Errorf("Incorrect method: Got %q, Want %q", output.Method, want)
	}
	wantParams := []Pair{}
	if !reflect.DeepEqual(wantParams, output.Params) {
		t.Errorf("Incorrect params: Got %v, Want %v", output.Params, wantParams)
	}
}

func TestPost_ParamsArg(t TB) {
	args := url.Values{
		"c": {"d"},
		"a": {"b", "g"},
	}
	got, err := gophernet.Post("http://localhost:"+PORT+"/echo", args)
	if err != nil {
		t.Fatalf("POST failed with error %v", err)
	}
	var output Request
	if err := json.Unmarshal(got, &output); err != nil {
		t.Fatalf("JSON Unmarshal failed with error %v", err)
	}
	if want := "POST"; output.Method != want {
		t.Errorf("Incorrect method: Got %q, Want %q", output.Method, want)
	}
	wantParams := []Pair{{"a", "b"}, {"a", "g"}, {"c", "d"}}
	if !reflect.DeepEqual(wantParams, output.Params) {
		t.Errorf("Incorrect params: Got %v, Want %v", output.Params, wantParams)
	}
}

func TestPost_ParamsMix(t TB) {
	args := url.Values{
		"c": {"d"},
		"a": {"b", "g"},
	}
	got, err := gophernet.Post("http://localhost:"+PORT+"/echo?a=l&m=n", args)
	if err != nil {
		t.Fatalf("POST failed with error %v", err)
	}
	var output Request
	if err := json.Unmarshal(got, &output); err != nil {
		t.Fatalf("JSON Unmarshal failed with error %v", err)
	}
	if want := "POST"; output.Method != want {
		t.Errorf("Incorrect method: Got %q, Want %q", output.Method, want)
	}
	wantParams := []Pair{{"a", "b"}, {"a", "g"}, {"c", "d"}}
	if !reflect.DeepEqual(wantParams, output.Params) {
		t.Errorf("Incorrect params: Got %v, Want %v", output.Params, wantParams)
	}
}
