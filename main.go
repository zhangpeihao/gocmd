package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

const (
	programName = "gocmd"
	version     = "0.1"
)

var (
	bindAddress *string = flag.String("BindAddress", ":8001", "The bind address.")
	root        *string = flag.String("Root", "/root/src", "The root path.")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s version[%s]\r\nUsage: %s [OPTIONS]\r\n", programName, version, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	http.HandleFunc("/", cmdHandler)

	//	http.HandleFunc("/deploy", deployHandler)
	//	http.HandleFunc("/set_tag", setTagHandler)
	//	http.HandleFunc("/start", startHandler)
	//	http.HandleFunc("/test", testHandler)
	go http.ListenAndServe(*bindAddress, nil)
	catchSignal()

}

func catchSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	os.Exit(0)
}

func cmdHandler(w http.ResponseWriter, req *http.Request) {
	url := strings.Split(strings.Split(req.RequestURI, "/")[1], "?")[0]
	params := strings.Split(req.FormValue("params"), ",")

	fmt.Printf("%s: %s(%v)\n", time.Now().Format(time.RFC3339), url, params)
	cmdString := path.Join(*root, url) + ".sh"
	cmd := exec.Command(cmdString, params...)
	response(w, cmd)
}

func response(w http.ResponseWriter, cmd *exec.Cmd) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(400)
		htmlOutput(w, fmt.Sprintf("Error: %s", err.Error()))
		return
	}
	w.WriteHeader(200)
	htmlOutput(w, string(out))
}

func htmlOutput(w io.Writer, str string) {
	b := []byte(str)
	b = bytes.Replace(b, []byte{27, 91, 48, 109}, []byte("</font>"), -1)
	b = bytes.Replace(b, []byte{27, 91, 52, 52, 109}, []byte(`<font color="blue" size="5">`), -1)
	b = bytes.Replace(b, []byte{27, 91, 51, 50, 109}, []byte(`<font color="green" size="6">`), -1)
	b = bytes.Replace(b, []byte{27, 91, 51, 49, 109}, []byte(`<font color="red" size="6">`), -1)
	b = bytes.Replace(b, []byte{10}, []byte(`<BR/>`), -1)
	w.Write([]byte(`<html>`))
	w.Write(b)
	w.Write([]byte(`</html>`))
}

