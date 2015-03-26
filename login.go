package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/secure"
)

func LoginForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  T("login", "", map[string]string{
    "return": r.URL.Query().Get("return"),
  })(w, r, ps)
}

func Login (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  uid := r.FormValue("uid")
  // TODO: validate password, maybe fetch roles
  // pwd := r.FormValue("pwd")
  secure.LogIn(w, r, uid)
  http.Redirect(w, r, r.FormValue("return"), 302)
}
