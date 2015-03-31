package main

import (
  r "github.com/dancannon/gorethink"
  "log"
)

type Database struct{
  Name string
  Session *r.Session
}

func InitDB (address, database string) *Database {
  db := Database{
    Name: database,
  }
  var err error
  db.Session, err = r.Connect(r.ConnectOpts {
    Address:  address,
    Database: database,
  })
  if err != nil {
    log.Fatalln(err.Error())
  }
  return &db
}
