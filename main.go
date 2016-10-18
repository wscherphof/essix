package main

import (
	"github.com/wscherphof/env"
	"os/exec"
	"os"
)

var (
	gopath = env.Get("GOPATH")
)

func main() {
	script := gopath + "/src/github.com/wscherphof/essix/script/essix"
	cmd := exec.Command(script, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
