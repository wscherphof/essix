package main

import (
  "github.com/wscherphof/msg"
)

func init () {
  var m, a = msg.Definition()

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

  m("Time format")
  a("nl", "2-1-2006, 15.04 u. (MST)")
  a("en", "1/2/2006 at 3:04pm (MST)")

  m("Last modified")
  a("nl", "Laatste wijziging")
  a("en", "Last modified")

  m("Save changes")
  a("nl", "Wijzigingen opslaan")
  a("en", "Save changes")

  m("Current")
  a("nl", "Huidig")
  a("en", "Current")

  m("New")
  a("nl", "Nieuw")
  a("en", "New")

  m("change")
  a("nl", "wijzigen")
  a("en", "change")

  m("Confirm")
  a("nl", "Bevestigen")
  a("en", "Confirm")
}
