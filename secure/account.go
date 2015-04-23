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
    "link": scheme + "://" + r.Host + "/account/activation/" + account.UID + "?code=" + account.ActivationCode,
    "name": account.Name(),
  })(r)
  return email.Send(subject, string(body), account.UID)
}

func SignUp (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  var (
    acc *account.Account
    err error
    conflict bool
  ) 
  if acc, err, conflict = account.New(r.FormValue); err != nil {
    if conflict {
      util.Error(w, r, ps, err, http.StatusConflict)
    } else {
      util.Error(w, r, ps, err)
    }
  } else if err = activationEmail(r, acc); err != nil && err != email.ErrNotSentImmediately {
    util.Error(w, r, ps, err)
  }
  if acc == nil {
    w.Write(util.BTemplate("signup_error-tail", "", nil)(r))
  } else {
    data := map[string]interface{}{
      "uid": acc.UID,
      "name": acc.Name(),
      "remark": "",
    }
    if err == email.ErrNotSentImmediately {
      data["remark"] = email.ErrNotSentImmediately.Error()
    }
    util.Template("signup_success", "", data)(w, r, ps)
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
    w.Write(util.BTemplate("activate_error-tail", "", nil)(r))
  } else {
    util.Template("activate_success", "", map[string]interface{}{
      "name": account.Name(),
    })(w, r, ps)
  }
}

func ActivationCodeForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  util.Template("activate_resend", "", map[string]interface{}{
    "uid": ps.ByName("uid"),
  })(w, r, ps)
}

func ActivationCode (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  if acc, err, conflict := account.GetInsecure(r.FormValue("uid")); err != nil {
    if conflict {
      util.Error(w, r, ps, err, http.StatusConflict)
    } else {
      util.Error(w, r, ps, err)
    }
    // TODO
    // w.Write(util.BTemplate("activate_error-tail", "", nil)(r))
  } else if acc.IsActive() {
      util.Error(w, r, ps, account.ErrAlreadyActivated, http.StatusConflict)
  } else if err := activationEmail(r, acc); err != nil && err != email.ErrNotSentImmediately {
    util.Error(w, r, ps, err)
  } else {
    // TODO
    // util.Template("activate_success", "", map[string]interface{}{
    //   "name": acc.Name(),
    // })(w, r, ps)
  }
}
