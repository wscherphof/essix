package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/expeertise/model/account"
  // "log"
)

func LogInForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  // TODO: captcha
  util.Template("login", "", nil)(w, r, ps)
}

func LogIn (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  handle := util.Handle(w, r, ps)
  if account, err, conflict := account.Get(r.FormValue("uid"), r.FormValue("pwd")); err != nil {
    handle(err, conflict, "login", nil)
  } else if err := secure.LogIn(w, r, account, true); err != nil {
    handle(err, false, "login", nil)
  }
}

func LogOut (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  secure.LogOut(w, r, true)
}
