package secure

import (
	"github.com/wscherphof/entity"
	"github.com/wscherphof/secure"
	"log"
)

type secureDB struct {
	*entity.Base
	*secure.Config
}

/*
Fetch implements github.com/wscherphof/secure.SecureDB.
*/
func (s *secureDB) Fetch(dst *secure.Config) (err error) {
	var empty bool
	if err, empty = s.Read(s); err != nil {
		if empty {
			log.Println("WARNING: SecureDB.Fetch():", err)
		} else {
			log.Println("ERROR: SecureDB.Fetch():", err)
		}
	} else {
		*dst = *s.Config
	}
	return
}

/*
Upsert implements github.com/wscherphof/secure.SecureDB.
*/
func (s *secureDB) Upsert(src *secure.Config) (err error) {
	s.Config = src
	if err = s.Update(s); err != nil {
		log.Println("ERROR: SecureDB.Upsert():", err)
	}
	return
}
