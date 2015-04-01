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

func NewPWD (pwd1, pwd2 string) (pwd PWD, err error) {
  if pwd1 != pwd2 {
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

func SignUp (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  if pwd, err := NewPWD(r.FormValue("pwd1"), r.FormValue("pwd2")); err != nil {
    Error(w, r, ps, err)
  } else if res, err := db.Insert(PWD_TABLE, pwd); err != nil {
    Error(w, r, ps, err)
  } else if len(res.GeneratedKeys) != 1 {
    Error(w, r, ps, errors.New("Failed saving password"))
  } else {
    // TODO: further validation, unique key on UID, password hashing, confirmation email, ...
    record := Account{
      Created: time.Now(),
      UID: r.FormValue("uid"),
      PWD: res.GeneratedKeys[0],
      Country: r.FormValue("country"),
      Postcode: r.FormValue("postcode"),
      FirstName: r.FormValue("firstname"),
      LastName: r.FormValue("lastname"),
    }
    if _, err := db.Insert(ACCOUNT_TABLE, record); err != nil {
      Error(w, r, ps, err)
    } else {
      w.WriteHeader(http.StatusOK)
      w.Write([]byte("account created"))
    }
  }
}
