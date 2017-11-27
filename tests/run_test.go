package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	mcl "github.com/njasm/marionette_client"
)

func handlerHome(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func handlerScript(w http.ResponseWriter, _ *http.Request, script string) {
	data, err := ioutil.ReadFile(script)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%s", data)
}

func handlerJSHome(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w,
		`<html><head><script src="/script.js"></script></head></html>`)
}

func handlerEcho(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	var query url.Values
	if r.Method == "POST" {
		query = r.PostForm
	} else {
		query = r.Form
	}
	params := []Pair{}
	for name, arr := range query {
		for _, val := range arr {
			params = append(params, Pair{name, val})
		}
	}
	req := Request{
		Method: r.Method,
		Params: sortPair(params),
	}

	out, err := json.Marshal(req)
	if err != nil {
		panic(err) // []Request is not cyclic, this will never be hit.
	}
	fmt.Fprintf(w, "%s", out)
}

func handlerBinary(w http.ResponseWriter, r *http.Request) {
	binary := deterministicBinData()
	fmt.Fprintf(w, "%s", binary)
}

func Init(t *testing.T) string {
	dir, err := ioutil.TempDir("", "golang")
	if err != nil {
		t.Fatal("Error creating temp directory:", err)
	}
	filename := filepath.Join(dir, "output.js")

	listener, err := net.Listen("tcp", "0.0.0.0:"+PORT)
	if err != nil {
		panic(err)
	}

	handlerFile := func(w http.ResponseWriter, r *http.Request) {
		handlerScript(w, r, filename)
	}
	http.HandleFunc("/", handlerHome)
	http.HandleFunc("/js", handlerJSHome)
	http.HandleFunc("/echo", handlerEcho)
	http.HandleFunc("/binary", handlerBinary)
	http.HandleFunc("/script.js", handlerFile)
	go http.Serve(listener, nil)

	return filename
}

func TestAll(t *testing.T) {
	script := Init(t)
	t.Run("TestAllJS", func(t *testing.T) { testAllJS(t, script) })
	t.Run("TestAllGo", testAllGo)
}

func testAllJS(t *testing.T, filename string) {
	t.Parallel()
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		t.Fatal("GOPATH variable is not set")
	}

	testDir := filepath.Join(
		gopath, "src/github.com/pallavagarwal07/gophernet/tests")
	fileList, err := ioutil.ReadDir(testDir)
	if err != nil {
		t.Fatal("ReadDir on test dir failed:", err)
	}

	files := []string{}
	for _, finfo := range fileList {
		if n := finfo.Name(); !strings.Contains(n, "_test.go") {
			files = append(files, finfo.Name())
		}
	}

	args := []string{"build", "-o", filename}
	buildOut, err := exec.Command("gopherjs", append(args, files...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("Gopherjs compilation failed: %v, %s", err, buildOut)
	}

	cmd := exec.Command("firefox", "--headless", "--marionette", "--disable-gpu")
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal("Getting pipe failed:", err)
	}

	err = cmd.Start()
	if err != nil {
		t.Fatal("Starting FF failed:", err)
	}
	output := []byte{}
	c := time.Tick(10 * time.Second)
	for {
		tmpOut := make([]byte, 100, 100)
		read, err := outPipe.Read(tmpOut)
		if err != nil {
			t.Fatal("Pipe read failed:", err, string(output))
		}
		if read != 0 {
			output = append(output, tmpOut...)
		}
		if bytes.Contains(output, []byte("Listening")) {
			break
		}
		select {
		case _ = <-c:
			break
		default:
			// try again
		}
	}

	if !bytes.Contains(output, []byte("Listening")) {
		t.Fatal("Could not get listening process")
	}
	runTest(filename, t)

	if err := cmd.Process.Kill(); err != nil {
		t.Fatal("Could not kill process:", err)
	}
}

func runTest(fname string, t *testing.T) {
	client := mcl.NewClient()
	if err := client.Connect("", 0); err != nil {
		t.Fatal("Client connect failed:", err)
	}
	_, err := client.NewSession("", nil)
	if err != nil {
		t.Fatal("New Session failed:", err)
	}
	client.Navigate("http://localhost:" + PORT + "/js")

	isElemValue := func() (string, error) {
		script := `
		if(typeof window.my_important_result !== 'undefined' && window.my_important_result[0] != '-') {
			return window.my_important_result;
		}
		return '--';`
		v, err := client.ExecuteScript(script, nil, 1000, false)
		if err != nil {
			return "", err
		}
		final := getVal(v.Value)
		if final == "--" {
			return "", errors.New("Not yet")
		}
		return final, nil
	}

	str, err := isElemValue()
	for ; err != nil; str, err = isElemValue() {
		time.Sleep(time.Millisecond * 50)
	}

	var results []Result
	if err := json.Unmarshal([]byte(str), &results); err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	for _, test := range results {
		t.Run(test.Name, func(t *testing.T) {
			if test.Fail {
				t.Fatal(test.Outp)
			}
		})
	}
}

func testAllGo(t *testing.T) {
	t.Parallel()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) { test(t) })
	}
}
