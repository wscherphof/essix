package main

import (
  "net/http"
  "errors"
  "github.com/julienschmidt/httprouter"
  // "log"
)

func SignUpForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  T("signup", "", map[string]interface{}{
    "Countries": Countries,
  })(w, r, ps)
}

func SignUp (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  // uid := r.FormValue("uid")
  pwd1 := r.FormValue("pwd1")
  pwd2 := r.FormValue("pwd2")
  if pwd1 != pwd2 {
    Error(w, r, ps, errors.New("Passwords not equal"))
  }
}
