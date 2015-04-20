package db

import (
  r "github.com/dancannon/gorethink"
  "log"
)

var (
  db string
  s *r.Session
)

func Init (address, database string) {
  db = database
  if session, err := r.Connect(r.ConnectOpts{
    Address:  address,
    Database: database,
  }); err != nil {
    log.Fatalln(err.Error())
  } else {
    s = session
  }
}

func NewInsertOpts () r.InsertOpts {
  return r.InsertOpts{}
}

func Insert (table string, record interface{}, opts ...r.InsertOpts) (r.WriteResponse, error) {
  return r.Table(table).Insert(record, opts...).RunWrite(s)
}

func Delete (table, key string) (r.WriteResponse, error) {
  return r.Table(table).Get(key).Delete().RunWrite(s)
}

func Truncate (table string) (r.WriteResponse, error) {
  return r.Table(table).Delete().RunWrite(s)
}

func NewTableCreateOpts () r.TableCreateOpts {
  return r.TableCreateOpts{}
}

func TableCreate (table string, opts ...r.TableCreateOpts) (*r.Cursor, error) {
  return r.Db(db).TableCreate(table, opts...).Run(s)
}

func Get (table, key string, result interface{}) (err error, found bool) {
  if cursor, e := r.Table(table).Get(key).Run(s); e != nil {
    err = e
  } else if e = cursor.One(result); e == nil {
    found = true
  } else if e != r.ErrEmptyResult {
    err = e
  }
  return
}

func One (table string, result interface{}) (err error, found bool) {
  if cursor, e := r.Table(table).Run(s); e != nil {
    err = e
  } else if e = cursor.One(result); e == nil {
    found = true
  } else if e != r.ErrEmptyResult {
    err = e
  }
  return
}

func All (table string) (cursor *r.Cursor, err error) {
  return r.Table(table).Run(s)
}
