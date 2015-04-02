package main

import (
  "net/http"
  "errors"
  "github.com/julienschmidt/httprouter"
  "time"
  "strings"
)

func SignUpForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  T("signup", "", map[string]interface{}{
    "Countries": Countries(),
  })(w, r, ps)
}

const PWD_TABLE = "pwd"
type PWD struct {
  Value string
}

func NewPWD (pwd1, pwd2 string) (pwd PWD, err error) {
  // TODO: further validation, password hashing, ...
  if pwd1 == "" {
    err = errors.New("Password empty")
  } else if pwd1 != pwd2 {
    err = errors.New("Passwords not equal")
  } else {
    pwd = PWD{
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

func NewAccount (val func (string) (string)) (account Account, err error, conflict bool) {
  uid := strings.ToLower(val("uid"))
  if got, _ := db.Get(ACCOUNT_TABLE, uid); got != nil {
    err, conflict = errors.New("Email address taken"), true
  } else if pwd, e := NewPWD(val("pwd1"), val("pwd2")); e != nil {
    err, conflict = e, true
  } else if res, e := db.Insert(PWD_TABLE, pwd); e != nil {
    err = e
  } else if len(res.GeneratedKeys) != 1 {
    err = errors.New("Unexpected error")
  } else {
    account = Account{
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
    }
  }
  return
}

func SignUp (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  if account, err, conflict := NewAccount(r.FormValue); err != nil {
    if conflict {
      Error(w, r, ps, err, http.StatusConflict)
    } else {
      Error(w, r, ps, err)
    }
  } else {
    // TODO: confirmation email, formatted response, ...
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("account created: " + account.UID))
  }
}
