package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/secure"
)

func LogOut (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  secure.LogOut(w, r)
}
