package secure

import (
  "github.com/wscherphof/msg"
  "github.com/wscherphof/secure"
)

func init () {
  var m, a = msg.Definition()

  m(secure.ErrNoTLS.Error())
  a("nl", "Voor inloggen is een vercijferde verbinding (https) vereist")
  a("en", "Logging in requires an encrypted connection (https)")

  m(secure.ErrTokenNotSaved.Error())
  a("nl", "Het lukt de webserver niet om de inlogsessie aan te maken :-(")
  a("en", "The server is failing to create the log in session :-(")

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

  m("Signup successful")
  a("nl", "Bedankt voor het aanmelden bij Expeertise! Voordat je kan inloggen moet je account nog geactiveerd worden; we hebben je een e-mail gestuurd met de activatiecode.")
  a("en", "Thanks for signing up with Expeertise! Before you can log in, your account needs to be activated; we've sent you an email containing the activation code.")

  m("activation subject")
  a("nl", "Je account activeren")
  a("en", "Activate your account")

  m("Activation code")
  a("nl", "Activeringscode")
  a("en", "Activation code")

  m("Activation code source")
  a("nl", "De activeringscode is je per e-mail toegestuurd")
  a("en", "The activation code was sent to you by email")

  m("Resend activation code")
  a("nl", "Stuur me de activeringscode opnieuw")
  a("en", "Resend me the activation code")

  m("Activate successful")
  a("nl", "Bedankt voor het activeren van je account bij Expeertise! Je registratie is compleet; je kan nu inloggen.")
  a("en", "Thanks for activating your account with Expeertise! Your registration is complete; you are now able to log in.")

  m("Resend successful")
  a("nl", "Dank voor je aanvraag. Controleer je e-mail voor de activatiecode.")
  a("en", "Thanks for your request. Check your email for the activation code.")

  m("password subject")
  a("nl", "Mijn wachtwoord opnieuw instellen")
  a("en", "Reset my password")

  m("Forgot password")
  a("nl", "Wachtwoord kwijt?")
  a("en", "Lost your password?")

  m("Passwordcode successful")
  a("nl", "Dank voor je aanvraag om je wachtwoord opnieuw in te stellen. Je krijgt een e-mail van ons met verdere instructies.")
  a("en", "Thanks for your request to reset your password. Please check your email for further instructions.")

  m("Passwordcode cancelled")
  a("nl", "Opnieuw instellen van het wachtwoord is geannuleerd.")
  a("en", "Password reset is cancelled.")

  m("Create new password")
  a("nl", "Maak een nieuw wachtwoord aan")
  a("en", "Create a new password")

  m(ErrPasswordCodeTimedOut.Error())
  a("nl", "De geldigheidstermijn van de aanvraag voor het opnieuw instellen van het wachtwoord is verstreken")
  a("en", "The request for a password reset is expired")

  m("Change password successful")
  a("nl", "Je nieuwe wachtwoord is nu actief.")
  a("en", "Your new password is activated.")

  m("Expires")
  a("nl", "Geldig tot")
  a("en", "Expires")

  m("Edit account")
  a("nl", "Mijn account")
  a("en", "My account details")

  m("Complete account")
  a("nl", "We hebben nog een paar laatste gegevens van je nodig")
  a("en", "Please provide the following concluding details")

  m("Star is required")
  a("nl", "Alleen de velden met een * zijn verplicht")
  a("en", "Only the fields with a * are required")
}
