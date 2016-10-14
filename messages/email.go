package messages

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.Key("Change email").
		Set("nl", "Je e-mailadres wijzigen").
		Set("en", "Changing your email address")

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

	msg.Key("Change email successful").
		Set("nl", "Je nieuwe e-mailadres is nu actief.").
		Set("en", "Your new email address is activated.")
}
