package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/expeertise/secure"
)

func Protected (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  account := secure.Authentication(r)
  util.Template("protected", "lang", map[string]interface{}{
    "name": account.Name(),
  })(w, r, ps)
}
