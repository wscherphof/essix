package secure

import (
  "github.com/wscherphof/msg"
  "github.com/wscherphof/secure"
)

func DefineMessages () {
  var m, a = msg.Init()

  m(secure.ErrTokenNotSaved.Error())
  a("nl", "Het lukt de webserver niet om de inlogsessie aan te maken :-(")
  a("en", "The server is failing to create the log in session :-(")
}
