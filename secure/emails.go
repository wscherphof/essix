package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/msg"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/wscherphof/expeertise/email"
)

func sendEmail (r *http.Request, acc *account.Account, resource, code string) (error) {
  subject := msg.Msg(r)(resource + " subject")
  scheme := "http"
  if r.TLS != nil {
    scheme = "https"
  }
  path := scheme + "://" + r.Host + "/account/" + resource + "/" + acc.UID
  // TODO: format links as "buttons" instead of URLs
  body := util.BTemplate(resource + "_email", "lang", map[string]interface{}{
    "action": path + "?code=" + code,
    "cancel": path,
    "name": acc.Name(),
  })(r)
  return email.Send(subject, string(body), acc.UID)
}

func activationEmail (r *http.Request, acc *account.Account) (error) {
  return sendEmail (r, acc, "activation", acc.ActivationCode)
}

func ActivationCodeForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  util.Template("activation_resend", "", map[string]interface{}{
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
    w.Write(util.BTemplate("activation_resend_error-tail", "", map[string]interface{}{
      "uid": r.FormValue("uid"),
    })(r))
  } else if acc.IsActive() {
      util.Error(w, r, ps, account.ErrAlreadyActivated, http.StatusConflict)
  } else if err := activationEmail(r, acc); err != nil && err != email.ErrNotSentImmediately {
    util.Error(w, r, ps, err)
  } else {
    util.Template("activation_resend_success", "", map[string]interface{}{
      "name": acc.Name(),
      "uid": acc.UID,
    })(w, r, ps)
  }
}

func passwordEmail (r *http.Request, acc *account.Account) (error) {
  // TODO: indicate Expires-time
  return sendEmail (r, acc, "password", acc.PasswordCode.Value)
}

func PasswordCodeForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  util.Template("passwordcode", "", map[string]interface{}{
    "uid": ps.ByName("uid"),
  })(w, r, ps)
}

func PasswordCode (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  if acc, err, conflict := account.GetInsecure(r.FormValue("uid")); err != nil {
    if conflict {
      util.Error(w, r, ps, err, http.StatusConflict)
    } else {
      util.Error(w, r, ps, err)
    }
    w.Write(util.BTemplate("passwordcode_error-tail", "", map[string]interface{}{
      "uid": r.FormValue("uid"),
    })(r))
  } else if err := acc.CreatePasswordCode(); err != nil {
      util.Error(w, r, ps, err)
  } else if err := passwordEmail(r, acc); err != nil && err != email.ErrNotSentImmediately {
    util.Error(w, r, ps, err)
  } else {
    util.Template("passwordcode_success", "", map[string]interface{}{
      "name": acc.Name(),
    })(w, r, ps)
  }
}
