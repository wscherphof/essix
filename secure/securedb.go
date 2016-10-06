package secure

import (
	"github.com/wscherphof/essix/entity"
	"github.com/wscherphof/secure"
	"log"
)

type config struct {
	*entity.Base
	*secure.Config
}

var conf = &config{
	Base: &entity.Base{
		ID: "secure",
	},
}

func init() {
	entity.Register(conf, "config")
}

type secureDB struct{}

func (s *secureDB) Fetch(dst *secure.Config) (err error) {
	if e := conf.Read(conf); e != nil {
		err = e
		log.Println("WARNING: SecureDB.Fetch():", err)
	} else {
		*dst = *conf.Config
	}
	return
}

func (s *secureDB) Upsert(src *secure.Config) (err error) {
	conf.Config = src
	if err = conf.Update(conf); err != nil {
		log.Println("WARNING: SecureDB.Upsert():", err)
	}
	return
}
