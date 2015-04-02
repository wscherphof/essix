package model

import (
  "github.com/wscherphof/expeertise/db"
  "errors"
  "time"
  "strings"
)

const PWD_TABLE = "pwd"
type PWD struct {
  Value string
}

func NewPWD (pwd1, pwd2 string) (pwd *PWD, err error) {
  // TODO: further validation, password hashing, ...
  if pwd1 == "" {
    err = errors.New("Password empty")
  } else if pwd1 != pwd2 {
    err = errors.New("Passwords not equal")
  } else {
    pwd = &PWD{
      Value: pwd1,
    }
  }
  return
}

const ACCOUNT_TABLE = "account"
type Account struct {
  Created time.Time
  UID string
  PWD string
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
    err, conflict = errors.New("Email address taken"), true
  } else if pwd, e := NewPWD(val("pwd1"), val("pwd2")); e != nil {
    err, conflict = e, true
  } else if res, e := db.Insert(PWD_TABLE, pwd); e != nil {
    err = e
  } else if len(res.GeneratedKeys) != 1 {
    err = errors.New("Unexpected error")
  } else {
    account = &Account{
      Created: time.Now(),
      UID: uid,
      PWD: res.GeneratedKeys[0],
      Country: val("country"),
      Postcode: strings.ToUpper(val("postcode")),
      FirstName: val("firstname"),
      LastName: val("lastname"),
    }
    if _, err = db.Insert(ACCOUNT_TABLE, account); err != nil {
      db.Delete(PWD_TABLE, account.PWD)
      account = nil
    }
  }
  return
}
