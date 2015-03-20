package main

import (
  "github.com/wscherphof/secure"
  "time"
)

type SecureDB struct {}

func (SecureDB) Fetch () secure.Config {
  return secure.Config {
    Secret: secure.Secret {
      Key: "qwerty",
      Time: time.Now(),
    },
    RedirectPath: "/login",
    TimeOut: 15 * 60 * time.Second, 
  }
}

func InitSecure () {
  secure.Init (SecureDB {})
}
