package db

import (
	"github.com/wscherphof/entity"
	"github.com/wscherphof/secure"
	"log"
)

type config struct {
	*entity.Base
	*secure.Config
}

func init() {
	entity.Register(&config{}, "config")
}

func New() *db {
	return &db{}
}

type db struct{}

var conf = &config{
	Base: &entity.Base{
		ID: "secure",
	},
}

func (*db) Fetch(dst *secure.Config) (err error) {
	var empty bool
	if err, empty = conf.Read(conf); err != nil {
		if empty {
			log.Println("WARNING: SecureDB.Fetch():", err)
		} else {
			log.Println("ERROR: SecureDB.Fetch():", err)
		}
	} else {
		*dst = *conf.Config
	}
	return
}

func (*db) Upsert(src *secure.Config) (err error) {
	conf.Config = src
	if err = conf.Update(conf); err != nil {
		log.Println("ERROR: SecureDB.Upsert():", err)
	}
	return
}
