package router

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.New(ErrInternalServerError.Error()).
		Add("nl", `Sorry; er is op de server iets onverwachts misgegaan.
			De foutmelding is gelogd voor onze systeembeheerders.
			Probeer het later opnieuw.`).
		Add("en", `Sorry; something unexpected went wrong on the server.
			De error is logged for our system operators.
			Please try again later.`)
}
