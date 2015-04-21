package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/data"
  "github.com/wscherphof/expeertise/model"
  "github.com/wscherphof/expeertise/email"
)

func SignUpForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  // TODO: captcha
  T("signup", "", map[string]interface{}{
    "Countries": data.Countries(),
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
    // TODO: formatted email, activation link & function, formatted response, ...
    if err := email.Send("Test email from Go!", "This is the email body.", account.UID); err != nil {
      Error(w, r, ps, err)
    } else if err == email.ErrNotSentImmediately {
      // 
    } else {
      w.WriteHeader(http.StatusOK)
      w.Write([]byte("account created: " + account.UID))
    }
  }
}
