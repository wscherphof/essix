package main

import (
  "github.com/wscherphof/secure"
  r "github.com/dancannon/gorethink"
  "log"
)

func InitSecure (db *Database) {
  secure.Init(newSecureDB(db))
}

type SecureDB struct {
  table string
  session *r.Session
}

func newSecureDB (db *Database) *SecureDB {
  s := &SecureDB {
    table: "secureConfig",
    session: db.Session,
  }
  // Create the table if needed
  if _, err := r.Table(s.table).Info().Run(s.session); err != nil {
    if _, err := r.Db(db.Name).TableCreate(s.table).Run(s.session); err != nil {
      log.Panicln("ERROR: ", err.Error())
    }
  }
  return s
}

func (s *SecureDB) Fetch () *secure.Config {
  if rows, err := r.Table(s.table).Run(s.session); err == nil {
    config := new(secure.Config)
    if err := rows.One(config); err == nil {
      return config
    }
  }
  return nil
}

func (s *SecureDB) Upsert (config *secure.Config) {
  if _, err := r.Table(s.table).Delete().RunWrite(s.session); err == nil {
    if _, err := r.Table(s.table).Insert(config).RunWrite(s.session); err != nil {
      log.Panicln("ERROR: ", err.Error())
    }
  }
}
