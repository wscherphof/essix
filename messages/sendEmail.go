package messages

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.Key("ErrNotSentImmediately").
		Set("nl", "We konden het e-mailbericht niet direct verzenden, maar je zou het spoedig alsnog moeten ontvangen.").
		Set("en", "We couldn't send the email message immediately, but you should receive it shortly")
}
