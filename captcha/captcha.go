package captcha

import (
	"github.com/dchest/captcha"
	"github.com/wscherphof/expeertise/db"
	"log"
	"time"
)

const (
	table   = "captcha"
	timeout = 15 * time.Minute
)

type captchaType struct {
	ID      string `gorethink:"id"`
	Digits  []byte
	Created int64
}

type store struct{}

func (s *store) Set(id string, digits []byte) {
	if _, err := db.Insert(table, &captchaType{
		ID:      id,
		Digits:  digits,
		Created: time.Now().Unix(),
	}); err != nil {
		log.Println("ERROR: Insert failed in table "+table+":", err)
	}
}

func (s *store) Get(id string, clear bool) (digits []byte) {
	c := new(captchaType)
	if err, found := db.Get(table, id, c); err != nil {
		log.Println("ERROR: Get failed in table "+table+":", err)
	} else if !found {
		log.Println("INFO: Not found in table "+table+":", id)
	} else {
		digits = c.Digits
	}
	return
}

var Server = captcha.Server(captcha.StdWidth, captcha.StdHeight)

func init() {
	if _, err := db.TableCreate(table); err == nil {
		log.Println("INFO: table created:", table)
		if _, err := db.IndexCreate(table, "Created"); err != nil {
			log.Println("ERROR: failed to create index:", table, err)
		} else {
			log.Println("INFO: index created:", table)
		}
	}
	captcha.SetCustomStore(new(store))
	go func() {
		for {
			limit := time.Now().Unix()
			time.Sleep(timeout)
			if _, err := db.DeleteTerm(db.Between(table, "Created", nil, limit, true, true)); err != nil {
				log.Println("WARNING: captcha prune failed:", err)
			}
		}
	}()
}
