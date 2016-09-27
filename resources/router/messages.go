package router

import (
	"github.com/wscherphof/msg"
	pkg "github.com/wscherphof/expeertise/router"
)

func Init() {
	msg.New(pkg.ErrInternalServerError.Error()).
		Add("nl", `Sorry; er is op de server iets onverwachts misgegaan.
			De foutmelding is gelogd voor onze systeembeheerders.
			Probeer het later opnieuw.`).
		Add("en", `Sorry; something unexpected went wrong on the server.
			De error is logged for our system operators.
			Please try again later.`)

	msg.New("Take me home").
		Add("nl", "Naar de startpagina").
		Add("en", "To the home page")

	msg.New("Try again").
		Add("nl", "Opnieuw proberen").
		Add("en", "Try again")
}
