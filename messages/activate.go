package messages

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.Key("Activate account").
		Set("nl", "Je account activeren").
		Set("en", "Activating your account")

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
}
