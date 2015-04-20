package config

import (
  "github.com/wscherphof/expeertise/db"
  "log"
  "errors"
)

const CONFIG_TABLE string = "config"

func Init () {
  opts := db.NewTableCreateOpts()
  opts.PrimaryKey = "Key"
  if cursor, _ := db.TableCreate(CONFIG_TABLE, opts); cursor != nil {
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
  _, err = db.Insert(CONFIG_TABLE, record)
  return
}
