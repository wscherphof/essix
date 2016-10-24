package main

import (
	"app/messages"
	"app/routes"
	essix "github.com/wscherphof/essix/server"
)

func init() {
	messages.Init()
	routes.Init()
}

func main() {
	essix.Run()
}
