package account

import (
  "github.com/wscherphof/msg"
)

func DefineMessages () {
  var m, a = msg.Init()

  m(ErrInvalidCredentials.Error())
  a("nl", "E-mailadres onbekend of wachtwoord of activatiecode onjuist")
  a("en", "Unknown email address or incorrect password or activation code")

  m(ErrPasswordEmpty.Error())
  a("nl", "Het wachtwoord mag niet leeg zijn")
  a("en", "Password must not be empty")

  m(ErrPasswordsNotEqual.Error())
  a("nl", "De ingevoerde wachtwoorden in beide velden moeten identiek zijn")
  a("en", "Both password entries must be identical")

  m(ErrEmailTaken.Error())
  a("nl", "Dit e-mailadres wordt al gebruikt voor een bestaand account")
  a("en", "This email address is already used for an existing account")

  m(ErrNotActivated.Error())
  a("nl", "Dit account moet eerst nog geactiveerd worden voordat je ermee kan inloggen")
  a("en", "This account still needs to be activated before it can be used to log in with")

  m(ErrAlreadyActivated.Error())
  a("nl", "Dit account is al geactiveerd")
  a("en", "This account is already activated")
}
