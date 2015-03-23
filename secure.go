package main

import (
  "github.com/wscherphof/secure"
  r "github.com/dancannon/gorethink"
  "log"
)

func InitSecure () {
  secure.Init(newSecureDB())
}

type SecureDB struct {
  session *r.Session
  db string
  table string
}

func newSecureDB () *SecureDB {
  s := &SecureDB {
    db: "expeertise",
    table: "secureConfig",
  }
  // TODO: store the session someplace higher
  if session, err := r.Connect(r.ConnectOpts {
    Address:  "localhost:28015",
    Database: s.db,
  }); err == nil {
    s.session = session
  } else {
    log.Fatalln(err.Error())
  }
  // Create the table if needed
  if _, err := r.Table(s.table).Info().Run(s.session); err != nil {
    if _, err := r.Db(s.db).TableCreate(s.table).Run(s.session); err != nil {
      log.Fatalln(err.Error())
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
  if _, err := r.Table(s.table).Delete().RunWrite(s.session); err != nil {
    log.Fatalln(err.Error())
  }
  if _, err := r.Table(s.table).Insert(config).RunWrite(s.session); err != nil {
    log.Fatalln(err.Error())
  }
}
