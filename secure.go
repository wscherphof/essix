package main

import (
  "github.com/wscherphof/secure"
  r "github.com/dancannon/gorethink"
  "log"
)

type SecureDB struct {
  Session *r.Session
}

func (s SecureDB) Fetch () *secure.Config {
  rows, err := r.Table("secureConfig").Nth(0).Run(s.Session)
  if err != nil {
    log.Fatalln(err.Error())
  }
  var config secure.Config
  err = rows.One(&config)
  if err != nil {
    log.Fatalln(err.Error())
  }
  return &config
}

func (s SecureDB) Update (config *secure.Config) {
  _, err := r.Table("secureConfig").Update(config).RunWrite(s.Session)
  if err != nil {
    log.Fatalln(err.Error())
  }
}

func InitSecure () {
  // TODO: store the session someplace higher
  session, err := r.Connect(r.ConnectOpts{
    Address:  "localhost:28015",
    Database: "expeertise",
  })
  if err != nil {
    log.Fatalln(err.Error())
  }
  secureDB := SecureDB {
    Session: session,
  }
  secure.Init(&secureDB)
}
