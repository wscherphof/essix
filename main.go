/*
Package main runs the `essix` command.

[Essix](https://github.com/wscherphof/essix) runs an essential simple secure
stable scalable stateless server.

	$ go get -u github.com/wscherphof/essix

The `essix` [command](https://github.com/wscherphof/essix#essix-command) manages
Essix apps, their server certificates (TLS/HTTPS), their backend databases
(RethinkDB), and their infrastructure (Docker Swarm Mode).

Follow the [Quickstart](https://github.com/wscherphof/essix#quickstart) to get
your first app running within minutes, on a swarm near you.

With `$ essix init`, a new Essix app package is inititialised, much like
https://github.com/wscherphof/essix/tree/master/app, where you would add your
own functionality. It includes the Profile example of how things can be done.

*/
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
