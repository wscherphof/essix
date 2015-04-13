package model

import (
  "github.com/wscherphof/expeertise/db"
  "errors"
  "time"
  "strings"
  "golang.org/x/crypto/bcrypt"
  "log"
)

var (
  ErrInvalidCredentials = errors.New("Unknown email address or incorrect password")
  ErrPasswordEmpty = errors.New("Password empty")
  ErrPasswordsNotEqual = errors.New("Passwords not equal")
  ErrEmailTaken = errors.New("Email address taken")
)

const ACCOUNT_TABLE = "account"

func InitAccount () {
  opts := db.NewTableCreateOpts()
  opts.PrimaryKey = "UID"
  if cursor, _ := db.TableCreate(ACCOUNT_TABLE, opts); cursor != nil {
    log.Println("INFO: table created:", ACCOUNT_TABLE)
  }
}

type Password struct {
 Created time.Time
 Value []byte
}

func NewPassword (pwd1, pwd2 string) (pwd *Password, err error) {
  if pwd1 == "" {
    err = ErrPasswordEmpty
  } else if pwd1 != pwd2 {
    err = ErrPasswordsNotEqual
  } else if hash, e := bcrypt.GenerateFromPassword([]byte(pwd1), bcrypt.DefaultCost); err != nil {
    err = e
  } else {
    pwd = &Password{
      Created: time.Now(),
      Value: hash,
    }
  }
  return
}

type Account struct {
  Created time.Time
  UID string
  PWD Password
  Country string
  Postcode string
  FirstName string
  LastName string
}

func NewAccount (val func (string) (string)) (account *Account, err error, conflict bool) {
  uid := strings.ToLower(val("uid"))
  if e, found := db.Get(ACCOUNT_TABLE, uid, new(Account)); e != nil {
    err = e
  } else if found {
    err, conflict = ErrEmailTaken, true
  } else if pwd, e := NewPassword(val("pwd1"), val("pwd2")); e != nil {
    err, conflict = e, true
  } else {
    account = &Account{
      Created: time.Now(),
      UID: uid,
      PWD: *pwd,
      Country: val("country"),
      Postcode: strings.ToUpper(val("postcode")),
      FirstName: val("firstname"),
      LastName: val("lastname"),
    }
    if _, err = db.Insert(ACCOUNT_TABLE, account); err != nil {
      account = nil
    }
  }
  return
}

func GetAccount (uid, pwd string) (account *Account, err error) {
  acc := new(Account)
  if e, found := db.Get(ACCOUNT_TABLE, strings.ToLower(uid), acc); e != nil {
    err = e
  } else if !(found) {
    err = ErrInvalidCredentials
  } else if e := bcrypt.CompareHashAndPassword(acc.PWD.Value, []byte(pwd)); e != nil {
    err = ErrInvalidCredentials
  } else {
    account = acc
  }
  return
}
