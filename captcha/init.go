package captcha

import (
  "github.com/dchest/captcha"
  "github.com/wscherphof/expeertise/db"
  "log"
  "time"
)

const CAPTCHA_TABLE string = "captcha"
const CAPTCHA_BUF int = 100
const CAPTCHA_TRIES int = 3
const TIMEOUT time.Duration = 15 * time.Minute

type captchaType struct{
  ID string `gorethink:"id"`
  Digits []byte
  Created int64
}

type store struct{}

var captchaChannel chan *captchaType

func setHandler () {
  for item := range captchaChannel {
    if _, err := db.Insert(CAPTCHA_TABLE, item); err != nil {
      log.Println("ERROR: Insert failed in table " + CAPTCHA_TABLE + ":", err)
    }
  }
}

func (s *store) Set (id string, digits []byte) {
  // Wonder if this will ever prove just a bit too tricky
  captchaChannel <- &captchaType{
    ID: id,
    Digits: digits,
    Created: time.Now().Unix(),
  }
}

func (s *store) Get (id string, clear bool) (digits []byte) {
  c := new(captchaType)
  for try := 0; try < CAPTCHA_TRIES; try++ {
    if err, found := db.Get(CAPTCHA_TABLE, id, c); err != nil {
      log.Println("ERROR: Get failed in table " + CAPTCHA_TABLE + ":", err)
    } else if ! found {
      log.Println("ERROR: Not found in table " + CAPTCHA_TABLE + ":", id)
    } else {
      digits = c.Digits
      break
    }
    time.Sleep(time.Duration(try + 1) * time.Second)
  }
  if len(digits) == 0 {
    log.Println("DEBUG: GET FAILED!!!")
  }
  return
}

var Server = captcha.Server(captcha.StdWidth, captcha.StdHeight)

func Init () {
  DefineMessages()
  if cursor, _ := db.TableCreate(CAPTCHA_TABLE); cursor != nil {
    log.Println("INFO: table created:", CAPTCHA_TABLE)
    if _, err := db.IndexCreate(CAPTCHA_TABLE, "Created"); err != nil {
      log.Println("ERROR: failed to create index:", CAPTCHA_TABLE, err)
    } else {
      log.Println("INFO: index created:", CAPTCHA_TABLE)
    }
  }
  captchaChannel = make(chan *captchaType, CAPTCHA_BUF)
  go setHandler()
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
