package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/msg"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/expeertise/data"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/wscherphof/expeertise/email"
)

func SignUpForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  // TODO: captcha
  util.Template("signup", "", map[string]interface{}{
    "Countries": data.Countries(),
  })(w, r, ps)
}

func activationEmail (r *http.Request, account *account.Account) (error) {
  subject := msg.Msg(r)("Activate account")
  scheme := "http"
  if r.TLS != nil {
    scheme = "https"
  }
  body := util.BTemplate("activate_email", "lang", map[string]interface{}{
    "link": scheme + "://" + r.Host + r.URL.Path + "/activation/" + account.UID + "?code=" + account.ActivationCode,
    "name": account.Name(),
  })(r)
  return email.Send(subject, string(body), account.UID)
}

func SignUp (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  if account, err, conflict := account.New(r.FormValue); err != nil {
    if conflict {
      util.Error(w, r, ps, err, http.StatusConflict)
    } else {
      util.Error(w, r, ps, err)
    }
  } else if err := activationEmail(r, account); err != nil && err != email.ErrNotSentImmediately {
    util.Error(w, r, ps, err)
  // TODO: formatted response
  } else if err == email.ErrNotSentImmediately {
  } else {
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("account created: " + account.UID))
  }
}

func ActivateForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  util.Template("activate", "", map[string]interface{}{
    "uid": ps.ByName("uid"),
    "code": r.URL.Query().Get("code"),
  })(w, r, ps)
}

func Activate (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  if account, err, conflict := account.Activate(r.FormValue("uid"), r.FormValue("code")); err != nil {
    if conflict {
      util.Error(w, r, ps, err, http.StatusConflict)
    } else {
      util.Error(w, r, ps, err)
    }
  } else {
    // TODO: formatted response
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("account activated: " + account.UID))
  }
}
