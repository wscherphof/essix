package main

import (
	"github.com/wscherphof/env"
	"os"
	"os/exec"
)

var (
	gopath = env.Get("GOPATH")
)

func main() {
	script := gopath + "/src/github.com/wscherphof/essix/script/essix"
	os.Args[0] = script
	cmd := exec.Command("bash", os.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
