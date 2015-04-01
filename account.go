package main

import (
  "net/http"
  "errors"
  "github.com/julienschmidt/httprouter"
  "time"
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

func NewPWD (r *http.Request) (pwd PWD, err error) {
  // TODO: further validation, password hashing, ...
  pwd1, pwd2 := r.FormValue("pwd1"), r.FormValue("pwd2")
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

func NewAccount (r *http.Request, pwd string) (account Account, err error) {
  // TODO: further validation, ...
  account = Account{
    Created: time.Now(),
    UID: r.FormValue("uid"),
    PWD: pwd,
    Country: r.FormValue("country"),
    Postcode: r.FormValue("postcode"),
    FirstName: r.FormValue("firstname"),
    LastName: r.FormValue("lastname"),
  }
  return
}

func SignUp (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  if pwd, err := NewPWD(r); err != nil {
    Error(w, r, ps, err)
  } else if res, err := db.Insert(PWD_TABLE, pwd); err != nil {
    Error(w, r, ps, err)
  } else if len(res.GeneratedKeys) != 1 {
    Error(w, r, ps, errors.New("Failed saving password"))
  } else if account, err := NewAccount(r, res.GeneratedKeys[0]); err != nil {
    Error(w, r, ps, err)
  } else if res, err := db.Insert(ACCOUNT_TABLE, account); err != nil {
    db.Delete(PWD_TABLE, account.PWD)
    Error(w, r, ps, err)
  } else if res.Errors > 0 {
    // Can't tell why this isn't returned in err :-(
    if res.FirstError[0:9] == "Duplicate" {
      err = errors.New("Duplicate primary key")
    } else {
      err = errors.New("Unexpected error")
    }
    db.Delete(PWD_TABLE, account.PWD)
    Error(w, r, ps, err)
  } else {
    // TODO: confirmation email, formatted response, ...
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("account created"))
  }
}
