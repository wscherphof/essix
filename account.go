package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/model"
)

func SignUpForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  T("signup", "", map[string]interface{}{
    "Countries": Countries(),
  })(w, r, ps)
}

func SignUp (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  if account, err, conflict := model.NewAccount(r.FormValue); err != nil {
    if conflict {
      Error(w, r, ps, err, http.StatusConflict)
    } else {
      Error(w, r, ps, err)
    }
  } else {
    // TODO: confirmation email, formatted response, ...
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("account created: " + account.UID))
  }
}
