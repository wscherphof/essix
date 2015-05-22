package captcha

import (
  "github.com/dchest/captcha"
  "github.com/wscherphof/expeertise/db"
  "log"
)

const CAPTCHA_TABLE string = "captcha"

var Server = captcha.Server(captcha.StdWidth, captcha.StdHeight)

type captchaType struct{
  ID string `gorethink:"id,omitempty"`
  Digits []byte
}

type store struct{}

func (s *store) Set (id string, digits []byte) {
  if _, err := db.Insert(CAPTCHA_TABLE, &captchaType{
    ID: id,
    Digits: digits,
  }); err != nil {
    log.Println("ERROR: Insert failed in table " + CAPTCHA_TABLE + ":", err)
  }
}

func (s *store) Get (id string, clear bool) (digits []byte) {
  c := new(captchaType)
  if err, found := db.Get(CAPTCHA_TABLE, id, c); err != nil {
    log.Println("ERROR: Get failed in table " + CAPTCHA_TABLE + ":", err)
  } else if ! found {
    log.Println("ERROR: Not found in table " + CAPTCHA_TABLE + ":", err)
  } else {
    digits = c.Digits
  }
  return
}

func Init () {
  if cursor, _ := db.TableCreate(CAPTCHA_TABLE); cursor != nil {
    log.Println("INFO: table created:", CAPTCHA_TABLE)
  }
  captcha.SetCustomStore(new(store))
}
