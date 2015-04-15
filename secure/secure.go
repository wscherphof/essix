package secure

import (
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/db"
  "github.com/wscherphof/expeertise/model"
  "log"
)

const SECURE_CONFIG_TABLE = "secureConfig"

func Init () {
  DefineMessages()
  if cursor, _ := db.TableCreate(SECURE_CONFIG_TABLE); cursor != nil {
    log.Println("INFO: SecureDB.Init() table created:", SECURE_CONFIG_TABLE)
  }
  secure.Init(model.Account{}, &secureDB{})
}

type secureDB struct {}

func (s *secureDB) Fetch () (config *secure.Config) {
  conf := new(secure.Config)
  if err, found := db.One(SECURE_CONFIG_TABLE, conf); err != nil {
    log.Println("ERROR: SecureDB.Fetch() failed with:", err)
  } else if found {
    config = conf    
  }
  return
}

func (s *secureDB) Upsert (config *secure.Config) {
  if _, err := db.Truncate(SECURE_CONFIG_TABLE); err != nil {
    log.Panicln("ERROR:", err)
  } else if _, err := db.Insert(SECURE_CONFIG_TABLE, config); err != nil {
    log.Panicln("ERROR:", err)
  }
}
