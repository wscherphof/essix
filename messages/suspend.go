package messages

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.Key("Suspend account").
		Set("nl", "Je account opzeggen").
		Set("en", "Suspending your account")

	msg.Key("Suspend sure").
		Set("nl", `Weet je zeker dat je je account bij Essix wil opzeggen?
			Je kan je later weliswaar weer opnieuw aanmelden, maar met het opzeggen gaan nu wel alle gekoppelde gegevens verloren.`).
		Set("en", `Are you sure you want to suspend your account?
			Though you can always sign up again later, all data currently linked to your account would get lost now.`)

	msg.Key("Yes, that's what I want").
		Set("nl", "Ja, dat wil ik").
		Set("en", "Yes, that's what I want")

	msg.Key("Suspend token successful").
		Set("nl", "Je hebt een aanvraag ingediend om je account op te zeggen. Je krijgt een e-mail van ons met verdere instructies.").
		Set("en", "You filed a request to suspend your account. For further instructions, please check your email.")

	msg.Key("Suspend token cancelled").
		Set("nl", "Opzeggen van het account is geannuleerd.").
		Set("en", "Termination of the account is cancelled.")

	msg.Key("Suspend successful").
		Set("nl", "Je account is nu verwijderd.").
		Set("en", "Your account has been deleted now.")
}
