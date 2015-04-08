package db

import (
  r "github.com/dancannon/gorethink"
  "log"
)

var Database string
var Session *r.Session

func Init (address, database string) {
  Database = database
  if s, err := r.Connect(r.ConnectOpts{
    Address:  address,
    Database: database,
  }); err != nil {
    log.Fatalln(err.Error())
  } else {
    Session = s
  }
}

func Insert (table string, record interface{}) (r.WriteResponse, error) {
  return r.Table(table).Insert(record).RunWrite(Session)
}

func Delete (table, key string) (r.WriteResponse, error) {
  return r.Table(table).Get(key).Delete().RunWrite(Session)
}

func Truncate (table string) (r.WriteResponse, error) {
  return r.Table(table).Delete().RunWrite(Session)
}

func TableCreate (table string) (*r.Cursor, error) {
  return r.Db(Database).TableCreate(table).Run(Session)
}

func Get (table, key string, result interface{}) (err error, found bool) {
  if cursor, e := r.Table(table).Get(key).Run(Session); e != nil {
    err = e
  } else if e = cursor.One(result); e == nil {
    found = true
  } else if e != r.ErrEmptyResult {
    err = e
  }
  return
}

func One (table string, result interface{}) (err error, found bool) {
  if cursor, e := r.Table(table).Run(Session); e != nil {
    err = e
  } else if e = cursor.One(result); e == nil {
    found = true
  } else if e != r.ErrEmptyResult {
    err = e
  }
  return
}
