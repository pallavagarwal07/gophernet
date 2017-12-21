package main

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strings"

	"github.com/pallavagarwal07/gophernet"
)

// We can't use the standard Go test toolchain because we need to run the same
// tests with phantomjs + gopherjs compiled js files. Renaming this file with
// *_test.go format causes the `go run` to refuse to compile the test_runner.go
// (used for the js tests).
var tests = map[string]func(t TB){
	"TestGet":                TestGet,
	"TestGet_Method":         TestGet_Method,
	"TestGet_Binary":         TestGet_Binary,
	"TestGet_Header":         TestGet_Header,
	"TestGet_ParamsURL":      TestGet_ParamsURL,
	"TestGet_ParamsArg":      TestGet_ParamsArg,
	"TestGet_ParamsMix":      TestGet_ParamsMix,
	"TestPostForm":           TestPostForm,
	"TestPostForm_Method":    TestPostForm_Method,
	"TestPostForm_Binary":    TestPostForm_Binary,
	"TestPostForm_ParamsURL": TestPostForm_ParamsURL,
	"TestPostForm_ParamsArg": TestPostForm_ParamsArg,
	"TestPostForm_ParamsMix": TestPostForm_ParamsMix,
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

func TestGet_Header(t TB) {
	c := gophernet.Client{
		Header: gophernet.Header{
			"H1": []string{"First", "Second"},
			"H2": []string{"Third"},
		},
	}
	got, err := c.Get("http://localhost:"+PORT+"/header", nil)
	if err != nil {
		t.Fatalf("GET failed with error %v", err)
	}
	var out gophernet.Header
	if err := json.Unmarshal(got, &out); err != nil {
		t.Fatalf("JSON Unmarshal failed with error %v", err)
	}
	for key, list := range c.Header {
		if got, want := strings.Join(out[key], ", "), strings.Join(list, ", "); got != want {
			t.Fatalf("Incorrect header %s. Got: %v, want: %v\n", key, got, want)
		}
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

func TestPostForm(t TB) {
	got, err := gophernet.PostForm("http://localhost:"+PORT, nil)
	if err != nil {
		t.Fatalf("POST failed with error %v", err)
	}
	if want := "Hello World!"; string(got) != want {
		t.Fatalf("Got output: %q, Want: %q", string(got), want)
	}
}

func TestPostForm_Binary(t TB) {
	got, err := gophernet.PostForm("http://localhost:"+PORT+"/binary", nil)
	if err != nil {
		t.Fatalf("POST failed with error %v", err)
	}
	want := deterministicBinData()
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Binary blob data does not match: Got %v, Want %v", got, want)
	}
}

func TestPostForm_Method(t TB) {
	got, err := gophernet.PostForm("http://localhost:"+PORT+"/echo", nil)
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

func TestPostForm_ParamsURL(t TB) {
	got, err := gophernet.PostForm("http://localhost:"+PORT+"/echo?a=b&c=d&a=g", nil)
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

func TestPostForm_ParamsArg(t TB) {
	args := url.Values{
		"c": {"d"},
		"a": {"b", "g"},
	}
	got, err := gophernet.PostForm("http://localhost:"+PORT+"/echo", args)
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

func TestPostForm_ParamsMix(t TB) {
	args := url.Values{
		"c": {"d"},
		"a": {"b", "g"},
	}
	got, err := gophernet.PostForm("http://localhost:"+PORT+"/echo?a=l&m=n", args)
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
