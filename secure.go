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

func (s *SecureDB) Fetch () *secure.Config {
  if rows, err := r.Table(s.Table).Run(db.Session); err == nil {
    config := new(secure.Config)
    if err := rows.One(config); err == nil {
      return config
    }
  }
  return nil
}

func (s *SecureDB) Upsert (config *secure.Config) {
  if _, err := r.Table(s.Table).Delete().RunWrite(db.Session); err == nil {
    if _, err := r.Table(s.Table).Insert(config).RunWrite(db.Session); err != nil {
      log.Panicln("ERROR:", err.Error())
    }
  }
}
