package email

import (
  "github.com/wscherphof/msg"
)

func DefineMessages () {
  var m, a = msg.Init()

  m(ErrNotSentImmediately.Error())
  a("nl", "Het e-mailbericht kon niet meteen worden verzonden; het wordt zo snel mogelijk alsnog verzonden")
  a("en", "The email message could not be sent immediately; it's queued to get sent shortly")
}
