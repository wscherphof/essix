package config

import (
	"errors"
	"github.com/wscherphof/expeertise/db"
	"log"
)

const (
	table = "config"
	pk    = "Key"
)

func init() {
	if cursor, _ := db.TableCreatePK(table, pk); cursor != nil {
		log.Println("INFO: table created:", table)
	}
}

func Get(key string, result interface{}) (err error) {
	if e, found := db.Get(table, key, result); e != nil {
		err = e
	} else if !(found) {
		err = errors.New("Key " + key + " not found in table " + table)
	}
	return
}

func Set(record interface{}) (err error) {
	_, err = db.InsertUpdate(table, record)
	return
}
