package messages

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.Key("ErrInvalidCredentials").
		Set("nl", "E-mailadres onbekend of wachtwoord of beveiligingscode onjuist").
		Set("en", "Unknown email address or incorrect password or security token")

	msg.Key("ErrPasswordEmpty").
		Set("nl", "Het wachtwoord mag niet leeg zijn").
		Set("en", "Password must not be empty")

	msg.Key("ErrPasswordsNotEqual").
		Set("nl", "De ingevoerde wachtwoorden in beide velden moeten identiek zijn").
		Set("en", "Both password entries must be identical")

	msg.Key("ErrEmailTaken").
		Set("nl", "Dit e-mailadres wordt al gebruikt voor een bestaand account").
		Set("en", "This email address is already used for an existing account")

	msg.Key("ErrNotActivated").
		Set("nl", "Dit account moet eerst nog geactiveerd worden voordat je ermee kan inloggen").
		Set("en", "This account still needs to be activated before it can be used to log in with")

	msg.Key("ErrAlreadyActivated").
		Set("nl", "Dit account is al geactiveerd").
		Set("en", "This account is already activated")

	msg.Key("ErrPasswordTokenTimedOut").
		Set("nl", "De geldigheidstermijn van de aanvraag voor het opnieuw instellen van het wachtwoord is verstreken").
		Set("en", "The request for a password reset is expired")
}
