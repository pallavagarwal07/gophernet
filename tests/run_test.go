package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	mcl "github.com/njasm/marionette_client"
)

var scriptName string

func handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Hello World!")
}

func handlerJS(w http.ResponseWriter, _ *http.Request) {
	data, err := ioutil.ReadFile(scriptName)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%s", data)
}

func handlerJShome(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, `<html><head><script src="/script.js"></script></head></html>`)
}

func Init() {
	listener, err := net.Listen("tcp", "0.0.0.0:"+PORT)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/js", handlerJShome)
	http.HandleFunc("/script.js", handlerJS)
	go http.Serve(listener, nil)
}

func TestAll(t *testing.T) {
	Init()
	t.Run("TestAllJS", testAllJS)
	t.Run("TestAllGo", testAllGo)
}

func testAllJS(t *testing.T) {
	cmdSetup := `
set -eu
tmpDir="$(mktemp -d)"
cd ${GOPATH}/src/github.com/pallavagarwal07/gophernet/tests
gopherjs build -o "${tmpDir}/output.js" $(ls -1 | grep -vP '_test.go')
echo ${tmpDir}/output.js
`
	jsFile, err := exec.Command("sh", "-c", cmdSetup).CombinedOutput()
	if err != nil {
		t.Fatal("Build failed:", err)
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
	runTest(strings.TrimSpace(string(jsFile)), t)

	if err := cmd.Process.Kill(); err != nil {
		t.Fatal("Could not kill process:", err)
	}
}

func runTest(fname string, t *testing.T) {
	scriptName = fname
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
	for name, test := range tests {
		t.Run(name, func(t *testing.T) { test(t) })
	}
}
