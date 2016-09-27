package secure

import (
	"github.com/wscherphof/msg"
	"github.com/wscherphof/secure"
)

func Init() {
	msg.New(secure.ErrNoTLS.Error()).
		Add("nl", "Voor inloggen is een vercijferde verbinding (https) vereist").
		Add("en", "Logging in requires an encrypted connection (https)")

	msg.New(secure.ErrTokenNotSaved.Error()).
		Add("nl", "Het lukt de webserver niet om de inlogsessie aan te maken :-(").
		Add("en", "The server is failing to create the log in session :-(")

	msg.New("Email").
		Add("nl", "E-mailadres").
		Add("en", "Email address")

	msg.New("Password").
		Add("nl", "Wachtwoord").
		Add("en", "Password")

	msg.New("Log in").
		Add("nl", "Inloggen").
		Add("en", "Log in")

	msg.New("Log out").
		Add("nl", "Uitloggen").
		Add("en", "Log out")

	msg.New("Sign up").
		Add("nl", "Een account aanmaken").
		Add("en", "Sign up")

	msg.New("Repeat password").
		Add("nl", "Wachtwoord nogmaals").
		Add("en", "Repeat password")

	msg.New("Passwords not equal").
		Add("nl", "De wachtwoorden zijn niet hetzelfde").
		Add("en", "Passwords don't match")

	msg.New("Country").
		Add("nl", "Land").
		Add("en", "Country")

	msg.New("Postcode").
		Add("nl", "Postcode").
		Add("en", "Postal code")

	msg.New("First name").
		Add("nl", "Voornaam").
		Add("en", "First name")

	msg.New("Last name").
		Add("nl", "Achternaam").
		Add("en", "Last name")

	msg.New("Signup successful").
		Add("nl", "Bedankt voor het aanmelden bij Expeertise! Voordat je kan inloggen moet je account nog geactiveerd worden; we hebben je een e-mail gestuurd met de activatiecode.").
		Add("en", "Thanks for signing up with Expeertise! Before you can log in, your account needs to be activated; we've sent you an email containing the activation code.")

	msg.New("activation subject").
		Add("nl", "Je account activeren").
		Add("en", "Activate your account")

	msg.New("Activation code").
		Add("nl", "Activeringscode").
		Add("en", "Activation code")

	msg.New("Activation code source").
		Add("nl", "De activeringscode is je per e-mail toegestuurd").
		Add("en", "The activation code was sent to you by email")

	msg.New("Resend activation code").
		Add("nl", "Stuur me de activeringscode opnieuw").
		Add("en", "Resend me the activation code")

	msg.New("Activate successful").
		Add("nl", "Bedankt voor het activeren van je account bij Expeertise! Je registratie is compleet; je kan nu inloggen.").
		Add("en", "Thanks for activating your account with Expeertise! Your registration is complete; you are now able to log in.")

	msg.New("Resend successful").
		Add("nl", "Dank voor je aanvraag. Controleer je e-mail voor de activatiecode.").
		Add("en", "Thanks for your request. Check your email for the activation code.")

	msg.New("password subject").
		Add("nl", "Mijn wachtwoord opnieuw instellen").
		Add("en", "Reset my password")

	msg.New("Forgot password").
		Add("nl", "Wachtwoord kwijt?").
		Add("en", "Lost your password?")

	msg.New("Password code successful").
		Add("nl", "Dank voor je aanvraag om je wachtwoord opnieuw in te stellen. Je krijgt een e-mail van ons met verdere instructies.").
		Add("en", "Thanks for your request to reset your password. Please check your email for further instructions.")

	msg.New("Password code cancelled").
		Add("nl", "Opnieuw instellen van het wachtwoord is geannuleerd.").
		Add("en", "Password reset is cancelled.")

	msg.New("Create new password").
		Add("nl", "Maak een nieuw wachtwoord aan").
		Add("en", "Create a new password")

	msg.New("Change password successful").
		Add("nl", "Je nieuwe wachtwoord is nu actief.").
		Add("en", "Your new password is activated.")

	msg.New("Expires").
		Add("nl", "Geldig tot").
		Add("en", "Expires")

	msg.New("emailaddress subject").
		Add("nl", "Mijn e-mailadres wijzigen").
		Add("en", "Change my email address")

	msg.New("Change email address").
		Add("nl", "Je e-mailadres wijzigen").
		Add("en", "Change your email address")

	msg.New("Want replace").
		Add("nl", "Ik wil").
		Add("en", "I want to replace")

	msg.New("Replace with").
		Add("nl", "vervangen door").
		Add("en", "with")

	msg.New("Email address code successful").
		Add("nl", "Dank voor je aanvraag om je e-mailadres te wijzigen. Je krijgt op het nieuwe adres een e-mail van ons met verdere instructies.").
		Add("en", "Thanks for your request to change your email address. For further instructions, please check your email on the new address.")

	msg.New("Email address code cancelled").
		Add("nl", "Wijzigen e-mailadres is geannuleerd.").
		Add("en", "Changing email address is cancelled.")

	msg.New("Change email address successful").
		Add("nl", "Je nieuwe e-mailadres is nu actief.").
		Add("en", "Your new email address is activated.")

	msg.New("Edit account").
		Add("nl", "Mijn account").
		Add("en", "My account details")

	msg.New("Complete account").
		Add("nl", "We hebben nog een paar laatste gegevens van je nodig").
		Add("en", "Please provide the following concluding details")

	msg.New("Star is required").
		Add("nl", "Alleen de velden met een * zijn verplicht").
		Add("en", "Only the fields with a * are required")

	msg.New("terminate subject").
		Add("nl", "Mijn account opzeggen").
		Add("en", "Terminate my account")

	msg.New("Terminate account").
		Add("nl", "Je account opzeggen").
		Add("en", "Terminate your account")

	msg.New("Terminate sure").
		Add("nl", `Weet je zeker dat je je account bij Expeertise wil opzeggen?
			Je kan je later weliswaar weer opnieuw aanmelden, maar met het opzeggen gaan nu wel alle gekoppelde gegevens verloren.`).
		Add("en", `Are you sure you want to terminate your account?
			Though you can always sign up again later, all data currently linked to your account would get lost now.`)

	msg.New("Yes, that's what I want").
		Add("nl", "Ja, dat wil ik").
		Add("en", "Yes, that's what I want")

	msg.New("Terminate code successful").
		Add("nl", "Je hebt een aanvraag ingediend om je account op te zeggen. Je krijgt een e-mail van ons met verdere instructies.").
		Add("en", "You filed a request to terminate your account. For further instructions, please check your email.")

	msg.New("Terminate code cancelled").
		Add("nl", "Opzeggen van het account is geannuleerd.").
		Add("en", "Termination of the account is cancelled.")

	msg.New("Terminate successful").
		Add("nl", "Je account is nu verwijderd.").
		Add("en", "Your accounst has been deleted now.")
}
