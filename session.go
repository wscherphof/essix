package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/model/account"
  // "log"
)

func LogInForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  // TODO: captcha
  T("login", "", nil)(w, r, ps)
}

func LogIn (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  if account, err, conflict := account.Get(r.FormValue("uid"), r.FormValue("pwd")); err != nil {
    if conflict {
      Error(w, r, ps, err, http.StatusConflict)
    } else {
      Error(w, r, ps, err)
    }
  } else if err := secure.LogIn(w, r, account); err != nil {
    Error(w, r, ps, err)
  }
  // Won't see this on successful secure.LogIn, but doesn't do any harm
  w.Write(TB("login_error-tail", "", nil)(r))
}

func LogOut (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  secure.LogOut(w, r)
}
