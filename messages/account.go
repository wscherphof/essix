package messages

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.Key("Sign up").
		Set("nl", "Een account aanmaken").
		Set("en", "Sign up")

	msg.Key("Repeat password").
		Set("nl", "Wachtwoord nogmaals").
		Set("en", "Repeat password")

	msg.Key("Passwords not equal").
		Set("nl", "De wachtwoorden zijn niet hetzelfde").
		Set("en", "Passwords don't match")

	msg.Key("Signup successful").
		Set("nl", "Bedankt voor het aanmelden bij Essix! Voordat je kan inloggen moet je account nog geactiveerd worden; we hebben je een e-mail gestuurd met de activatiecode.").
		Set("en", "Thanks for signing up with Essix! Before you can log in, your account needs to be activated; we've sent you an email containing the activation code.")

	msg.Key("Edit account").
		Set("nl", "Mijn account").
		Set("en", "My account details")
}
