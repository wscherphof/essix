package secure

import (
	"github.com/wscherphof/essix/config"
	"github.com/wscherphof/secure"
	"log"
)

type secureConfigStore struct {
	Key   string
	Value *secure.Config
}

var store = &secureConfigStore{
	Key: "secure",
}

type secureDB struct{}

func (s *secureDB) Fetch(dst *secure.Config) (err error) {
	if err = config.Get(store.Key, store); err != nil {
		log.Println("WARNING: SecureDB.Fetch():", err)
	} else {
		*dst = *store.Value
	}
	return
}

func (s *secureDB) Upsert(src *secure.Config) (err error) {
	store.Value = src
	if err = config.Set(store); err != nil {
		log.Println("WARNING: SecureDB.Upsert():", err)
	}
	return
}
