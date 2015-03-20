package main

import (
  "github.com/wscherphof/msg"
)

func DefineMessages () {
  var m, a = msg.Init()
  m("")
  a("nl", "")
  a("en", "")

  m("User name")
  a("nl", "Gebruikersnaam")
  a("en", "User name")

  m("Password")
  a("nl", "Wachtwoord")
  a("en", "Password")

  m("Log in")
  a("nl", "Inloggen")
  a("en", "Log in")
}
