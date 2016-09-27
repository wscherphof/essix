package email

import (
	"github.com/wscherphof/msg"
	pkg "github.com/wscherphof/expeertise/email"
)

func Init() {
	msg.New(pkg.ErrNotSentImmediately.Error()).
		Add("nl", "We konden het e-mailbericht niet direct verzenden, maar je zou het spoedig alsnog moeten ontvangen.").
		Add("en", "We couldn't send the email message immediately, but you should receive it shortly")
}
