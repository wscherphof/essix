package account

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.New(ErrInvalidCredentials.Error()).
		Add("nl", "E-mailadres onbekend of wachtwoord of activeringscode onjuist").
		Add("en", "Unknown email address or incorrect password or activation code")

	msg.New(ErrPasswordEmpty.Error()).
		Add("nl", "Het wachtwoord mag niet leeg zijn").
		Add("en", "Password must not be empty")

	msg.New(ErrPasswordsNotEqual.Error()).
		Add("nl", "De ingevoerde wachtwoorden in beide velden moeten identiek zijn").
		Add("en", "Both password entries must be identical")

	msg.New(ErrEmailTaken.Error()).
		Add("nl", "Dit e-mailadres wordt al gebruikt voor een bestaand account").
		Add("en", "This email address is already used for an existing account")

	msg.New(ErrNotActivated.Error()).
		Add("nl", "Dit account moet eerst nog geactiveerd worden voordat je ermee kan inloggen").
		Add("en", "This account still needs to be activated before it can be used to log in with")

	msg.New(ErrAlreadyActivated.Error()).
		Add("nl", "Dit account is al geactiveerd").
		Add("en", "This account is already activated")

	msg.New(ErrPasswordCodeUnset.Error()).
		Add("nl", "Voor dit account is geen aanvraag bekend voor het opnieuw instellen van het wachtwoord").
		Add("en", "This account has no pending request for a password reset")

	msg.New(ErrPasswordCodeIncorrect.Error()).
		Add("nl", "Ongeldige aanvraag").
		Add("en", "Invalid request")

	msg.New(ErrEmailAddressCodeUnset.Error()).
		Add("nl", "Voor dit account is geen aanvraag bekend voor het wijzigen van het e-mailadres").
		Add("en", "This account has no pending request for changing the email address")

	msg.New(ErrEmailAddressCodeIncorrect.Error()).
		Add("nl", "Ongeldige aanvraag").
		Add("en", "Invalid request")
}
