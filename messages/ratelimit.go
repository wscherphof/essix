package messages

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.Key("ErrTooManyRequests").
		Set("nl", `Sorry; hier geldt een frequentielimiet; je kan dit maar eens in
      de zoveel tijd aanvragen.
      Probeer het later opnieuw.`).
		Set("en", `Sorry; a rate limit is in effect for this request type.
      Please try again later.`)

	msg.Key("ErrInvalidRequest").
		Set("nl", `Sorry; er klopt iets niet in het kader van de frequentielimiet
      die voor dit verzoek van kracht is.`).
		Set("en", `Sorry; something is wrong with the rate limit controls.`)

	msg.Key("Limit is set to").
		Set("nl", `De limiet is ingesteld op:`).
		Set("en", `Rate limit is set to:`)
}
