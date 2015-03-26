package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/secure"
)

func LogInForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  T("login", "", map[string]string{
    "return": r.URL.Query().Get("return"),
  })(w, r, ps)
}

func LogIn (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  uid := r.FormValue("uid")
  // TODO: validate password, maybe fetch roles
  // pwd := r.FormValue("pwd")
  secure.LogIn(w, r, uid)
}
