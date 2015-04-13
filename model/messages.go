package model

import (
  "github.com/wscherphof/msg"
)

func DefineMessages () {
  var m, a = msg.Init()

  m(ErrInvalidCredentials.Error())
  a("nl", "E-mailadres onbekend of wachtwoord onjuist")
  a("en", "Unknown email address or incorrect password")

  m(ErrPasswordEmpty.Error())
  a("nl", "Het wachtwoord mag niet leeg zijn")
  a("en", "Password must not be empty")

  m(ErrPasswordsNotEqual.Error())
  a("nl", "De ingevoerde wachtwoorden in beide velden moeten identiek zijn")
  a("en", "Both password entries must be identical")

  m(ErrEmailTaken.Error())
  a("nl", "Dit e-mailadres wordt al gebruikt voor een bestaand account")
  a("en", "This email address is already used for an existing account")
}
