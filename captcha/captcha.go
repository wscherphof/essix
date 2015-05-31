package captcha

import (
  "github.com/dchest/captcha"
  "github.com/wscherphof/expeertise/db"
  "log"
  "time"
)

const CAPTCHA_TABLE string = "captcha"
const TIMEOUT time.Duration = 15 * time.Minute

type captchaType struct{
  ID string `gorethink:"id"`
  Digits []byte
  Created int64
}

type store struct{}

func (s *store) Set (id string, digits []byte) {
  if _, err := db.Insert(CAPTCHA_TABLE, &captchaType{
    ID: id,
    Digits: digits,
    Created: time.Now().Unix(),
  }); err != nil {
    log.Println("ERROR: Insert failed in table " + CAPTCHA_TABLE + ":", err)
  }
}

func (s *store) Get (id string, clear bool) (digits []byte) {
  c := new(captchaType)
  if err, found := db.Get(CAPTCHA_TABLE, id, c); err != nil {
    log.Println("ERROR: Get failed in table " + CAPTCHA_TABLE + ":", err)
  } else if ! found {
    log.Println("INFO: Not found in table " + CAPTCHA_TABLE + ":", id)
  } else {
    digits = c.Digits
  }
  return
}

var Server = captcha.Server(captcha.StdWidth, captcha.StdHeight)

func init () {
  if cursor, _ := db.TableCreate(CAPTCHA_TABLE); cursor != nil {
    log.Println("INFO: table created:", CAPTCHA_TABLE)
    if _, err := db.IndexCreate(CAPTCHA_TABLE, "Created"); err != nil {
      log.Println("ERROR: failed to create index:", CAPTCHA_TABLE, err)
    } else {
      log.Println("INFO: index created:", CAPTCHA_TABLE)
    }
  }
  captcha.SetCustomStore(new(store))
  go func() {
    for {
      limit := time.Now().Unix()
      time.Sleep(TIMEOUT)
      if _, err := db.DeleteTerm(db.Between(CAPTCHA_TABLE, "Created", nil, limit, true, true)); err != nil {
        log.Println("WARNING: captcha prune failed:", err)
      }
    }
  }()
}
