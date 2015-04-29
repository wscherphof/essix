package main

import (
  "github.com/wscherphof/msg"
)

func DefineMessages () {
  var m, a = msg.Init()

  m("")
  a("nl", "")
  a("en", "")

  m("Hi")
  a("nl", "Hoi")
  a("en", "Hi")

  m("Take me home")
  a("nl", "Naar de startpagina")
  a("en", "To the home page")

  m("Try again")
  a("nl", "Opnieuw proberen")
  a("en", "Try again")
}
