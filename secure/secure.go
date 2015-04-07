package secure

import (
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/db"
  // TODO: ditch r dependency
  r "github.com/dancannon/gorethink"
  "log"
)

func Init () {
  const table = "secureConfig"
  // TODO: db.TableCreate
  // Create the table if needed
  if _, err := r.Table(table).Info().Run(db.Session); err != nil {
    if _, err := r.Db(db.Database).TableCreate(table).Run(db.Session); err != nil {
      log.Panicln("ERROR:", err.Error())
    }
  }
  secure.Init(&secureDB{
    table: table,
  })
}

type secureDB struct {
  table string
}

func (s *secureDB) Fetch () (config *secure.Config) {
  conf := new(secure.Config)
  if err, found := db.One(s.table, conf); err != nil {
    log.Println("ERROR: SecureDB.Fetch() failed with:", err.Error())
  } else if found {
    config = conf    
  } else {
    config = nil
  }
  return
}

func (s *secureDB) Upsert (config *secure.Config) {
  // TODO: db.Delete, db.Insert
  if _, err := r.Table(s.table).Delete().RunWrite(db.Session); err == nil {
    if _, err := r.Table(s.table).Insert(config).RunWrite(db.Session); err != nil {
      log.Panicln("ERROR:", err.Error())
    }
  }
}
