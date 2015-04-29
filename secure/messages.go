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

  m("Resend activation code")
  a("nl", "Stuur me de activeringscode opnieuw")
  a("en", "Resend me the activation code")

  m("Resend successful")
  a("nl", "Dank voor je aanvraag. Controleer je e-mail voor de activatiecode.")
  a("en", "Thanks for your request. Check your email for the activation code.")

  m("Email")
  a("nl", "E-mailadres")
  a("en", "Email address")

  m("Password")
  a("nl", "Wachtwoord")
  a("en", "Password")

  m("Log in")
  a("nl", "Inloggen")
  a("en", "Log in")

  m("Log out")
  a("nl", "Uitloggen")
  a("en", "Log out")

  m("Sign up")
  a("nl", "Een account aanmaken")
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
  a("nl", "Achternaam")
  a("en", "Last name")

  m("Activate account")
  a("nl", "Je account activeren")
  a("en", "Activate your account")

  m("Activation code")
  a("nl", "Activeringscode")
  a("en", "Activation code")

  m("Activation code source")
  a("nl", "De activeringscode is je per e-mail toegestuurd")
  a("en", "The activation code was sent to you by email")

  m("Activate")
  a("nl", "Activeren")
  a("en", "Activate")

  m("Try again")
  a("nl", "Opnieuw proberen")
  a("en", "Try again")
}
