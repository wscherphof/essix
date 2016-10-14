package messages

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.Key("Forgot password").
		Set("nl", "Wachtwoord kwijt?").
		Set("en", "Lost your password?")

	msg.Key("Reset password").
		Set("nl", "Je wachtwoord opnieuw instellen").
		Set("en", "Resetting your password")

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
}
