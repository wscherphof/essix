package secure

import (
	"github.com/wscherphof/expeertise/config"
	"github.com/wscherphof/secure"
	"log"
)

type secureDB struct{}

type secureConfigStore struct {
	Key   string
	Value *secure.Config
}

const key = "secure"

func (s *secureDB) Fetch() (conf *secure.Config) {
	store := new(secureConfigStore)
	if err := config.Get(key, store); err != nil {
		log.Println("WARNING: SecureDB.Fetch() failed with:", err)
	} else {
		conf = store.Value
	}
	return
}

func (s *secureDB) Upsert(conf *secure.Config) {
	if err := config.Set(&secureConfigStore{
		Key:   key,
		Value: conf,
	}); err != nil {
		log.Panicln("ERROR: SecureDB.Upsert():", err)
	}
}
