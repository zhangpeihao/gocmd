package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
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

	cmd := exec.Command(cmdString, mode)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	err = cmd.Start()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("Run err: %s", err.Error())))
		return
	}
	cmd.Wait()
	w.WriteHeader(200)

	n, err := io.Copy(w, stderr)
	m, err := io.Copy(w, stdout)
	if m+n == 0 {
		w.Write([]byte("no response"))
	}
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

	cmd := exec.Command(cmdString, mode)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	err = cmd.Start()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("Run err: %s", err.Error())))
		return
	}
	cmd.Wait()
	w.WriteHeader(200)

	n, err := io.Copy(w, stderr)
	m, err := io.Copy(w, stdout)
	if m+n == 0 {
		w.Write([]byte("no response"))
	}
}

func startHandler(w http.ResponseWriter, req *http.Request) {
	mode := req.FormValue("mode")
	module := req.FormValue("module")
	fmt.Printf("start(mode: %s, module: %s)\n", mode, module)
	if len(mode) == 0 || len(module) == 0 {
		w.WriteHeader(400)
		return
	}

	var cmdString string
	cmdString = path.Join(*root, "docker", "jk-docker", "start.sh")
	cmd := exec.Command(cmdString, mode)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	err = cmd.Start()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("Run err: %s", err.Error())))
		return
	}
	cmd.Wait()
	w.WriteHeader(200)

	n, err := io.Copy(w, stderr)
	m, err := io.Copy(w, stdout)
	if m+n == 0 {
		w.Write([]byte("no response"))
	}
}

