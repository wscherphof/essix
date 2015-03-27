package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/secure"
  // "log"
)

func LogInForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  T("login", "", nil)(w, r, ps)
}

func LogIn (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  uid := r.FormValue("uid")
  // TODO: validate password, maybe fetch roles
  // pwd := r.FormValue("pwd")
  // log.Print("DEBUG: TLS ", r.TLS)
  Error(w, r, ps, secure.LogIn(w, r, uid))
}

func LogOut (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  secure.LogOut(w, r)
}
