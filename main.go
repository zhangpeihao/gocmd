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

	http.HandleFunc("/deploy", deployHandler)
	http.HandleFunc("/set_tag", setTagHandler)
	http.HandleFunc("/start", startHandler)
	go http.ListenAndServe(*bindAddress, nil)
	catchSignal()

}

func catchSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
	os.Exit(0)
}

func deployHandler(w http.ResponseWriter, req *http.Request) {
	mode := req.FormValue("mode")
	module := req.FormValue("module")
	fmt.Printf("deploy(mode: %s, module: %s)\n", mode, module)
	if len(mode) == 0 || len(module) == 0 {
		w.WriteHeader(400)
		return
	}

	var cmdString string
	switch module {
	case "tomcat":
		cmdString = path.Join(*root, "java", "jk-mobile", "make.sh")
	case "h5":
		cmdString = path.Join(*root, "h5", "jk", "make.sh")
	case "h5-admin":
		cmdString = path.Join(*root, "h5", "jk-admin", "make.sh")
	case "static":
		cmdString = path.Join(*root, "h5", "static", "make.sh")
	default:
		w.WriteHeader(400)
		w.Write([]byte("Unknown module"))
		return
	}
	fmt.Println(cmdString, mode)

	cmd := exec.Command(cmdString, mode)
	response(w, cmd)
}

func setTagHandler(w http.ResponseWriter, req *http.Request) {
	mode := req.FormValue("mode")
	module := req.FormValue("module")
	fmt.Printf("set_tag(mode: %s, module: %s)\n", mode, module)
	if len(mode) == 0 || len(module) == 0 {
		w.WriteHeader(400)
		return
	}

	var cmdString string
	switch module {
	case "tomcat":
		cmdString = path.Join(*root, "java", "jk-mobile", "set_tag.sh")
	case "h5":
		cmdString = path.Join(*root, "h5", "jk", "set_tag.sh")
	case "h5-admin":
		cmdString = path.Join(*root, "h5", "jk-admin", "set_tag.sh")
	case "static":
		cmdString = path.Join(*root, "h5", "static", "set_tag.sh")
	default:
		w.WriteHeader(400)
		w.Write([]byte("Unknown module"))
		return
	}
	fmt.Println(cmdString, mode)

	cmd := exec.Command(cmdString, mode)
	response(w, cmd)
}

func startHandler(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	mode := req.FormValue("mode")
	module := req.FormValue("module")
	fmt.Printf("start(mode: %s, module: %s)\n", mode, module)
	if len(mode) == 0 || len(module) == 0 {
		w.WriteHeader(400)
		return
	}

	var cmdString string
	cmdString = path.Join(*root, "docker", "jk-docker", "start.sh")
	fmt.Println(cmdString, mode, module)

	cmd := exec.Command(cmdString, mode, module)
	response(w, cmd)
}

func response(w http.ResponseWriter, cmd *exec.Cmd) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(400)
		htmlOutput(w, fmt.Sprintf("<font color=\"red\">Error: %s</font><BR/>", err.Error()))
		if out != nil {
			htmlOutput(w, string(out))
		}
		return
	}
	w.WriteHeader(200)
	htmlOutput(w, string(out))
}

func htmlOutput(w io.Writer, str string) {
	str = strings.Replace(str, "[44m", `<font color="blue">`, -1)
	str = strings.Replace(str, "[32m", `<font color="green">`, -1)
	str = strings.Replace(str, "[31m", `<font color="red">`, -1)
	str = strings.Replace(str, "\n", "<BR/>", -1)
	str = string(bytes.Replace([]byte(str), []byte{91, 27, 48, 109}, []byte("</font>"), -1))
	w.Write([]byte(`<html>`))
	w.Write([]byte(str))
	w.Write([]byte(`</html>`))
}
