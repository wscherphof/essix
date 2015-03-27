package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/secure"
)

func Protected (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  uid := secure.Authenticate(w, r)
  T("protected", "lang", map[string]string{
    "uid": uid,
  })(w, r, ps)
}
