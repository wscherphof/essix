package account

import (
	"github.com/wscherphof/msg"
)

func Init() {
	msg.New("errInvalidCredentials").
		Add("nl", "E-mailadres onbekend of wachtwoord of activeringscode onjuist").
		Add("en", "Unknown email address or incorrect password or activation code")

	msg.New("errPasswordEmpty").
		Add("nl", "Het wachtwoord mag niet leeg zijn").
		Add("en", "Password must not be empty")

	msg.New("errPasswordsNotEqual").
		Add("nl", "De ingevoerde wachtwoorden in beide velden moeten identiek zijn").
		Add("en", "Both password entries must be identical")

	msg.New("errEmailTaken").
		Add("nl", "Dit e-mailadres wordt al gebruikt voor een bestaand account").
		Add("en", "This email address is already used for an existing account")

	msg.New("errNotActivated").
		Add("nl", "Dit account moet eerst nog geactiveerd worden voordat je ermee kan inloggen").
		Add("en", "This account still needs to be activated before it can be used to log in with")

	msg.New("errAlreadyActivated").
		Add("nl", "Dit account is al geactiveerd").
		Add("en", "This account is already activated")

	msg.New("errCodeUnset").
		Add("nl", "Geen aanvraag bekend voor dit account").
		Add("en", "No pending request for this account")

	msg.New("errCodeIncorrect").
		Add("nl", "Ongeldige aanvraag").
		Add("en", "Invalid request")

	msg.New("errPasswordCodeTimedOut").
		Add("nl", "De geldigheidstermijn van de aanvraag voor het opnieuw instellen van het wachtwoord is verstreken").
		Add("en", "The request for a password reset is expired")
}
