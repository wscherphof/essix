package main

import (
	essix "github.com/wscherphof/essix/server"
	"<messages>"
	"<routes>"
)

func init() {
	messages.Init()
	routes.Init()
}

func main() {
	essix.Run()
}
