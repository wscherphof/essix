package main

import (
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/db"
  r "github.com/dancannon/gorethink"
  "log"
)

func InitSecure () {
  secure.Init(newSecureDB())
}

type SecureDB struct {
  Table string
}

func newSecureDB () *SecureDB {
  const table = "secureConfig"
  // Create the table if needed
  if _, err := r.Table(table).Info().Run(db.Session); err != nil {
    if _, err := r.Db(db.Database).TableCreate(table).Run(db.Session); err != nil {
      log.Panicln("ERROR:", err.Error())
    }
  }
  return &SecureDB {
    Table: table,
  }
}

func (s *SecureDB) Fetch () (config *secure.Config) {
  if err, _ := db.One(s.Table, config); err != nil {
    log.Println("ERROR: SecureDB.Fetch() failed with:", err.Error())
  }
  return
}

func (s *SecureDB) Upsert (config *secure.Config) {
  if _, err := r.Table(s.Table).Delete().RunWrite(db.Session); err == nil {
    if _, err := r.Table(s.Table).Insert(config).RunWrite(db.Session); err != nil {
      log.Panicln("ERROR:", err.Error())
    }
  }
}
