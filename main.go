package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type Request struct {
	Cmd   string
	Os    string
	Stdin string
}

type Response struct {
	Stdout string
	Stderr string
}

func ExecCommand(req *Request) *Response {
	args := strings.Fields(req.Cmd)
	cmd := exec.Command(args[0], args[1:]...)
	stdin, _ := cmd.StdinPipe()
	io.WriteString(stdin, req.Stdin)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	errRun := cmd.Start()
	var err []byte
	var out []byte
	if errRun != nil {
		log.Println(errRun) // dublicate error. is_ok?
		err = []byte(fmt.Sprint(errRun))
	} else {
		out, _ = io.ReadAll(stdout)
		err, _ = io.ReadAll(stderr)
	}
	cmd.Wait()
	return &Response{
		Stdout: string(out),
		Stderr: string(err),
	}
}

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Empty request")
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	splitedRequest := strings.Split(string(body), "\n")
	ch_in := make(chan *Response, len(splitedRequest))
	for _, msg := range splitedRequest {
		if len(msg) > 3 && msg[1] == '{' && msg[len(msg)-2] == '}' { // 1/len(msg) - 2 - min left/right position of request
			var bodyRequest Request
			err := json.Unmarshal([]byte(msg[1:len(msg)-1]), &bodyRequest)
			if err != nil {
				log.Println(err)
				http.Error(w, "bad body", http.StatusBadRequest)
				return
			}

			go func(ch chan *Response, bodyRequest *Request) { // async execute commands
				ch <- ExecCommand(bodyRequest)
			}(ch_in, &bodyRequest)

		} else {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
	}

	ans := make([]*Response, len(splitedRequest))
	for i := 0; i < len(splitedRequest); i++ {
		ans[i] = <-ch_in
	}

	for _, elem := range ans {
		w.Write([]byte(fmt.Sprintf(`{ "stdout": %s, "stderr": %s}`+"\n", elem.Stdout, elem.Stderr)))
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	s := &http.Server{
		Addr:    ":8085",
		Handler: nil,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cer},
		},
	}

	go http.HandleFunc("/api/v1/remote-execution", HandleConnection) // already async ?
	s.ListenAndServeTLS("server.crt", "server.key")
}
