package secure

import (
  "github.com/wscherphof/msg"
  "github.com/wscherphof/secure"
)

func DefineMessages () {
  var m, a = msg.Init()

  m(secure.ErrTokenNotSaved.Error())
  a("nl", "Het lukt de webserver niet om de inlogsessie aan te maken :-(")
  a("en", "The server is failing to create the log in session :-(")

  m("Signup successful")
  a("nl", "Bedankt voor het aanmelden bij Expeertise! Voordat je kan inloggen moet je account nog geactiveerd worden; we hebben je een e-mail gestuurd met de activatiecode.")
  a("en", "Thanks for signing up with Expeertise! Before you can log in, your account needs to be activated; we've sent you an email containing the activation code.")

  m("Activate successful")
  a("nl", "Bedankt voor het activeren van je account bij Expeertise! Je registratie is compleet; je kan nu inloggen.")
  a("en", "Thanks for activating your account with Expeertise! Your registration is complete; you are now able to log in.")
}
