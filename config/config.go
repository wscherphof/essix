package config

import (
  "github.com/wscherphof/expeertise/db"
  "log"
  "errors"
)

const CONFIG_TABLE = "config"
const CONFIG_PK = "Key"

func init () {
  if cursor, _ := db.TableCreatePK(CONFIG_TABLE, CONFIG_PK); cursor != nil {
    log.Println("INFO: table created:", CONFIG_TABLE)
  }
}

func Get (key string, result interface{}) (err error) {
  if e, found := db.Get(CONFIG_TABLE, key, result); e != nil {
    err = e
  } else if !(found) {
    err = errors.New("Key " + key + " not found in table " + CONFIG_TABLE)
  }
  return
}

func Set (record interface{}) (err error) {
  _, err = db.InsertUpdate(CONFIG_TABLE, record)
  return
}
