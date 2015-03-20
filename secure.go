package main

import (
  "github.com/wscherphof/secure"
  "time"
)

func InitSecure () {
  secure.Init (func () secure.Config {
    // TODO: interact with a real DB ;-)
    return secure.Config {
      Secret: secure.Secret {
        Key: "qwerty",
        Time: time.Now(),
      },
      RedirectPath: "/login",
      TimeOut: 15 * 60 * time.Second, 
    }
  })
}
