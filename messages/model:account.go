package messages

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.New("ErrInvalidCredentials").
		Add("nl", "E-mailadres onbekend of wachtwoord of activeringscode onjuist").
		Add("en", "Unknown email address or incorrect password or activation code")

	msg.New("ErrPasswordEmpty").
		Add("nl", "Het wachtwoord mag niet leeg zijn").
		Add("en", "Password must not be empty")

	msg.New("ErrPasswordsNotEqual").
		Add("nl", "De ingevoerde wachtwoorden in beide velden moeten identiek zijn").
		Add("en", "Both password entries must be identical")

	msg.New("ErrEmailTaken").
		Add("nl", "Dit e-mailadres wordt al gebruikt voor een bestaand account").
		Add("en", "This email address is already used for an existing account")

	msg.New("ErrNotActivated").
		Add("nl", "Dit account moet eerst nog geactiveerd worden voordat je ermee kan inloggen").
		Add("en", "This account still needs to be activated before it can be used to log in with")

	msg.New("ErrAlreadyActivated").
		Add("nl", "Dit account is al geactiveerd").
		Add("en", "This account is already activated")

	msg.New("ErrCodeUnset").
		Add("nl", "Geen aanvraag bekend voor dit account").
		Add("en", "No pending request for this account")

	msg.New("ErrCodeIncorrect").
		Add("nl", "Ongeldige aanvraag").
		Add("en", "Invalid request")

	msg.New("ErrPasswordCodeTimedOut").
		Add("nl", "De geldigheidstermijn van de aanvraag voor het opnieuw instellen van het wachtwoord is verstreken").
		Add("en", "The request for a password reset is expired")
}
