package secure

import (
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/db"
  "log"
)

func Init () {
  const table = "secureConfig"
  if cursor, _ := db.TableCreate(table); cursor != nil {
    log.Println("INFO: SecureDB.Init() table created:", table)
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
    log.Println("ERROR: SecureDB.Fetch() failed with:", err)
  } else if found {
    config = conf    
  }
  return
}

func (s *secureDB) Upsert (config *secure.Config) {
  if _, err := db.Truncate(s.table); err != nil {
    log.Panicln("ERROR:", err)
  } else if _, err := db.Insert(s.table, config); err != nil {
    log.Panicln("ERROR:", err)
  }
}
