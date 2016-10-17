package main

import (
	essix "github.com/wscherphof/essix/server"
	"github.com/wscherphof/s6app/messages"
)

func init() {
	messages.Init()
}

func main() {
	essix.Run()
}
