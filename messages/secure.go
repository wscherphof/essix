package messages

import (
	"github.com/wscherphof/msg"
	"github.com/wscherphof/secure"
)

func init() {
	msg.Key(secure.ErrNoTLS.Error()).
		Set("nl", "Voor inloggen is een vercijferde verbinding (https) vereist").
		Set("en", "Logging in requires an encrypted connection (https)")

	msg.Key(secure.ErrTokenNotSaved.Error()).
		Set("nl", "Het lukt de webserver niet om de inlogsessie aan te maken :-(").
		Set("en", "The server is failing to create the log in session :-(")

	msg.Key("Email").
		Set("nl", "E-mailadres").
		Set("en", "Email address")

	msg.Key("Password").
		Set("nl", "Wachtwoord").
		Set("en", "Password")

	msg.Key("Log in").
		Set("nl", "Inloggen").
		Set("en", "Log in")

	msg.Key("Log out").
		Set("nl", "Uitloggen").
		Set("en", "Log out")

	msg.Key("Sign up").
		Set("nl", "Een account aanmaken").
		Set("en", "Sign up")

	msg.Key("Repeat password").
		Set("nl", "Wachtwoord nogmaals").
		Set("en", "Repeat password")

	msg.Key("Passwords not equal").
		Set("nl", "De wachtwoorden zijn niet hetzelfde").
		Set("en", "Passwords don't match")

	msg.Key("Country").
		Set("nl", "Land").
		Set("en", "Country")

	msg.Key("Postcode").
		Set("nl", "Postcode").
		Set("en", "Postal code")

	msg.Key("First name").
		Set("nl", "Voornaam").
		Set("en", "First name")

	msg.Key("Last name").
		Set("nl", "Achternaam").
		Set("en", "Last name")

	msg.Key("Signup successful").
		Set("nl", "Bedankt voor het aanmelden bij Essix! Voordat je kan inloggen moet je account nog geactiveerd worden; we hebben je een e-mail gestuurd met de activatiecode.").
		Set("en", "Thanks for signing up with Essix! Before you can log in, your account needs to be activated; we've sent you an email containing the activation code.")

	msg.Key("ActivateToken subject").
		Set("nl", "Je account activeren").
		Set("en", "Activate your account")

	msg.Key("Activate token").
		Set("nl", "Activeringscode").
		Set("en", "Activation token")

	msg.Key("Activate token source").
		Set("nl", "De activeringscode is je per e-mail toegestuurd").
		Set("en", "The activation token was sent to you by email")

	msg.Key("Resend activate token").
		Set("nl", "Stuur me de activeringscode opnieuw").
		Set("en", "Resend me the activation token")

	msg.Key("Activate successful").
		Set("nl", "Bedankt voor het activeren van je account bij Essix! Je registratie is compleet; je kan nu inloggen.").
		Set("en", "Thanks for activating your account with Essix! Your registration is complete; you are now able to log in.")

	msg.Key("Resend successful").
		Set("nl", "Dank voor je aanvraag. Controleer je e-mail voor de activatiecode.").
		Set("en", "Thanks for your request. Check your email for the activation token.")

	msg.Key("PasswordToken subject").
		Set("nl", "Mijn wachtwoord opnieuw instellen").
		Set("en", "Reset my password")

	msg.Key("Forgot password").
		Set("nl", "Wachtwoord kwijt?").
		Set("en", "Lost your password?")

	msg.Key("Password token successful").
		Set("nl", "Dank voor je aanvraag om je wachtwoord opnieuw in te stellen. Je krijgt een e-mail van ons met verdere instructies.").
		Set("en", "Thanks for your request to reset your password. Please check your email for further instructions.")

	msg.Key("Password token cancelled").
		Set("nl", "Opnieuw instellen van het wachtwoord is geannuleerd.").
		Set("en", "Password reset is cancelled.")

	msg.Key("Create new password").
		Set("nl", "Maak een nieuw wachtwoord aan").
		Set("en", "Create a new password")

	msg.Key("Change password successful").
		Set("nl", "Je nieuwe wachtwoord is nu actief.").
		Set("en", "Your new password is activated.")

	msg.Key("Expires").
		Set("nl", "Geldig tot").
		Set("en", "Expires")

	msg.Key("EmailToken subject").
		Set("nl", "Mijn e-mailadres wijzigen").
		Set("en", "Change my email address")

	msg.Key("Change email").
		Set("nl", "Je e-mailadres wijzigen").
		Set("en", "Change your email address")

	msg.Key("Want replace").
		Set("nl", "Ik wil").
		Set("en", "I want to replace")

	msg.Key("Replace with").
		Set("nl", "vervangen door").
		Set("en", "with")

	msg.Key("Email token successful").
		Set("nl", "Dank voor je aanvraag om je e-mailadres te wijzigen. Je krijgt op het nieuwe adres een e-mail van ons met verdere instructies.").
		Set("en", "Thanks for your request to change your email address. For further instructions, please check your email on the new address.")

	msg.Key("Email token cancelled").
		Set("nl", "Wijzigen e-mailadres is geannuleerd.").
		Set("en", "Changing email address is cancelled.")

	msg.Key("Change email  successful").
		Set("nl", "Je nieuwe e-mailadres is nu actief.").
		Set("en", "Your new email address is activated.")

	msg.Key("Edit account").
		Set("nl", "Mijn account").
		Set("en", "My account details")

	msg.Key("Complete account").
		Set("nl", "We hebben nog een paar laatste gegevens van je nodig").
		Set("en", "Please provide the following concluding details")

	msg.Key("Star is required").
		Set("nl", "Alleen de velden met een * zijn verplicht").
		Set("en", "Only the fields with a * are required")

	msg.Key("terminate subject").
		Set("nl", "Mijn account opzeggen").
		Set("en", "Terminate my account")

	msg.Key("Terminate account").
		Set("nl", "Je account opzeggen").
		Set("en", "Terminate your account")

	msg.Key("Terminate sure").
		Set("nl", `Weet je zeker dat je je account bij Essix wil opzeggen?
			Je kan je later weliswaar weer opnieuw aanmelden, maar met het opzeggen gaan nu wel alle gekoppelde gegevens verloren.`).
		Set("en", `Are you sure you want to terminate your account?
			Though you can always sign up again later, all data currently linked to your account would get lost now.`)

	msg.Key("Yes, that's what I want").
		Set("nl", "Ja, dat wil ik").
		Set("en", "Yes, that's what I want")

	msg.Key("Terminate token successful").
		Set("nl", "Je hebt een aanvraag ingediend om je account op te zeggen. Je krijgt een e-mail van ons met verdere instructies.").
		Set("en", "You filed a request to terminate your account. For further instructions, please check your email.")

	msg.Key("Terminate token cancelled").
		Set("nl", "Opzeggen van het account is geannuleerd.").
		Set("en", "Termination of the account is cancelled.")

	msg.Key("Terminate successful").
		Set("nl", "Je account is nu verwijderd.").
		Set("en", "Your account has been deleted now.")
}
