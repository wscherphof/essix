package main

import (
	essix "github.com/wscherphof/essix/server"
	"<app_package>/messages"
)

func init() {
	messages.Init()
}

func main() {
	essix.Run()
}
