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
  opts := r.InsertOpts{
    Conflict: "error",
    ReturnChanges: false,
  }
  return r.Table(table).Insert(record, opts).RunWrite(d.Session)
}

func (d *Database) Delete (table, key string) (r.WriteResponse, error) {
  return r.Table(table).Get(key).Delete().RunWrite(d.Session)
}

func (d *Database) Get (table, key string) (result interface{}, err error) {
  if cursor, err := r.Table(table).Get(key).Run(d.Session); err == nil {
    cursor.One(&result)
  }
  return
}
