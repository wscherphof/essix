package messages

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.Key("ErrInternalServerError").
		Set("nl", `Sorry; er is op de server iets onverwachts misgegaan.
			De foutmelding is gelogd voor onze systeembeheerders.
			Probeer het later opnieuw.`).
		Set("en", `Sorry; something unexpected went wrong on the server.
			De error is logged for our system operators.
			Please try again later.`)

	msg.Key("Take me home").
		Set("nl", "Naar de startpagina").
		Set("en", "To the home page")

	msg.Key("Try again").
		Set("nl", "Opnieuw proberen").
		Set("en", "Try again")

	msg.Key("ErrNotSentImmediately").
		Set("nl", "We konden het e-mailbericht niet direct verzenden, maar je zou het spoedig alsnog moeten ontvangen.").
		Set("en", "We couldn't send the email message immediately, but you should receive it shortly")
}
