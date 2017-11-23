package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
	"time"

	mcl "github.com/njasm/marionette_client"
)

func TestJS(t *testing.T) {
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

	cmd := exec.Command("setsid", "firefox", "--headless", "--marionette", "--disable-gpu")
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

func getVal(out string) string {
	var v interface{}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		return "--"
	}
	return v.(map[string]interface{})["value"].(string)
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
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatal("File read failed:", err)
	}
	script := string(data)
	_, err = client.ExecuteScript(script, nil, 1000, false)
	if err != nil {
		t.Fatal("Execute script failed:", err)
	}

	isElemValue := func() (string, error) {
		script := "if($('#output_test').length > 0) return $('#output_test').text(); return '--';"
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
		time.Sleep(time.Millisecond * 20)
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
