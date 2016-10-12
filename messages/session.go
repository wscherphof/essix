package messages

import (
	"github.com/wscherphof/msg"
	"github.com/wscherphof/secure"
)

func init() {
	msg.Key(secure.ErrTokenNotSaved.Error()).
		Set("nl", "Het lukt de webserver niet om de inlogsessie aan te maken :-(").
		Set("en", "The server is failing to create the log in session :-(")

	msg.Key("Email").
		Set("nl", "E-mailadres").
		Set("en", "Email address")

	msg.Key("Password").
		Set("nl", "Wachtwoord").
		Set("en", "Password")

	msg.Key("Forgot password").
		Set("nl", "Wachtwoord kwijt?").
		Set("en", "Lost your password?")

	msg.Key("Log in").
		Set("nl", "Inloggen").
		Set("en", "Log in")

	msg.Key("Log out").
		Set("nl", "Uitloggen").
		Set("en", "Log out")
}
