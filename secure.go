package main

import (
  "github.com/wscherphof/secure"
  r "github.com/dancannon/gorethink"
  "log"
)

func InitSecure () {
  // TODO: store the session someplace higher
  session, err := r.Connect(r.ConnectOpts{
    Address:  "localhost:28015",
    Database: "expeertise",
  })
  if err != nil {
    log.Fatalln(err.Error())
  }

  secure.Init (func () *secure.Config {
    rows, err := r.Table("secureConfig").Nth(0).Run(session)
    if err != nil {
      log.Fatalln(err.Error())
    }
    var conf secure.Config
    err = rows.One(&conf)
    if err != nil {
      log.Fatalln(err.Error())
    }
    return &conf
  })
}
