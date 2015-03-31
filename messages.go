package main

import (
  "github.com/wscherphof/msg"
  "github.com/wscherphof/secure"
)

func DefineMessages () {
  var m, a = msg.Init()
  m("")
  a("nl", "")
  a("en", "")

  m("Email")
  a("nl", "E-mailadres")
  a("en", "Email address")

  m("Password")
  a("nl", "Wachtwoord")
  a("en", "Password")

  m("Log in")
  a("nl", "Inloggen")
  a("en", "Log in")

  m(secure.ErrTokenNotSaved.Error())
  a("nl", "Het lukt de webserver niet om de inlogsessie aan te maken :-(")
  a("en", "The server is failing to create the log in session :-(")

  m("Log out")
  a("nl", "Uitloggen")
  a("en", "Log out")

  m("Sign up")
  a("nl", "Aanmelden")
  a("en", "Sign up")

  m("Repeat password")
  a("nl", "Wachtwoord nogmaals")
  a("en", "Repeat password")

  m("Passwords not equal")
  a("nl", "De wachtwoorden zijn niet hetzelfde")
  a("en", "Passwords don't match")

  m("Country")
  a("nl", "Land")
  a("en", "Country")

  m("Postcode")
  a("nl", "Postcode")
  a("en", "Postal code")

  m("First name")
  a("nl", "Voornaam")
  a("en", "First name")

  m("Last name")
  a("nl", "Achterrnaam")
  a("en", "Last name")
}
