package secure

import (
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/config"
  "github.com/wscherphof/expeertise/model/account"
  "log"
)

func Init () {
  DefineMessages()
  secure.Init(account.Account{}, &secureDB{})
}

type secureDB struct {}

type secureConfigStore struct{
  Key string
  Value *secure.Config
}

const SECURE_KEY = "secure"

func (s *secureDB) Fetch () (conf *secure.Config) {
  store := new(secureConfigStore)
  if err := config.Get(SECURE_KEY, store); err != nil {
    log.Println("WARNING: SecureDB.Fetch() failed with:", err)
  } else {
    conf = store.Value    
  }
  return
}

func (s *secureDB) Upsert (conf *secure.Config) {
  if err := config.Set(&secureConfigStore{
    Key: SECURE_KEY,
    Value: conf,
  }); err != nil {
    log.Panicln("ERROR: SecureDB.Upsert():", err)
  }
}
