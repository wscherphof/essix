package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/model"
  // "log"
)

func LogInForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  T("login", "", nil)(w, r, ps)
}

func LogIn (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  // TODO: maybe fetch roles
  // log.Print("DEBUG: TLS ", r.TLS)
  uid, pwd := r.FormValue("uid"), r.FormValue("pwd")
  if account, err := model.GetAccount(uid, pwd); err != nil {
    Error(w, r, ps, err)
  } else {
    Error(w, r, ps, secure.LogIn(w, r, account))
  }
}

func LogOut (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  secure.LogOut(w, r)
}
