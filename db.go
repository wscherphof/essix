package main

import (
  r "github.com/dancannon/gorethink"
  "log"
)

type Database struct{
  Name string
  Session *r.Session
}

func InitDB (address, database string) (ret *Database) {
  ret = &Database{
    Name: database,
  }
  var err error
  ret.Session, err = r.Connect(r.ConnectOpts {
    Address:  address,
    Database: database,
  })
  if err != nil {
    log.Fatalln(err.Error())
  }
  return
}

func (d *Database) Insert (table string, record interface{}) (r.WriteResponse, error) {
  return r.Table(table).Insert(record).RunWrite(d.Session)
}
