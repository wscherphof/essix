package email

import (
  "github.com/wscherphof/msg"
)

func init () {
  var m, a = msg.Init()

  m(ErrNotSentImmediately.Error())
  a("nl", "We konden het e-mailbericht niet direct verzenden, maar je zou het spoedig alsnog moeten ontvangen.")
  a("en", "We couldn't send the email message immediately, but you should receive it shortly")
}
