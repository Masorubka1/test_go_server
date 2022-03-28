package main

import (
	"log"
	"net/http"
	"bytes"
	"testing"
)

func TestExecCommand(t *testing.T) {
	in := Request{
		Cmd: "echo 'hello' 2>&1",
		Os: "linux",
		Stdin: "",
	}
	out := &Response {
		Stdout: "'hello' 2>&1\n",
		Stderr: "",
	}
	got := ExecCommand(&in)
	if got.Stdout != out.Stdout || got.Stderr != out.Stderr {
		t.Errorf("got %q, wanted %q", got, out)
	}
}

func TestHandleConnection(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	client := &http.Client{}
	url := "https://localhost:8085/api/v1/remote-execution"
	req := [5]string{``, `hello`, `[{"cmd": "df -h", "os": "linux", "stdin": ""}]`, `[{"cmd": "dddddd -g", "os": "linux", "stdin": ""}]`, `[{"cmd": "echo 'hello' 2>&1", "os": "linux", "stdin": ""}]`}
	ans := [5]string{"400 Bad Request", "400 Bad Request", "200 OK", "200 OK", "200 OK"}
	for ind, elem := range req {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(elem)))
		req.Header.Set("X-Custom-Header", "myvalue")
    	req.Header.Set("Content-Type", "application/json")
		
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.Status != ans[ind] {
			t.Errorf("got status %q, wanted %q", resp.Status, ans[ind])
		}
	}
}
