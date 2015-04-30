package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/msg"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/wscherphof/expeertise/email"
)

func sendEmail (r *http.Request, acc *account.Account, resource, code string) (err error, remark string) {
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
  if e := email.Send(subject, string(body), acc.UID); e != nil {
    if e == email.ErrNotSentImmediately {
      remark = e.Error()
    } else {
      err = e
    }
  }
  return
}

func activationEmail (r *http.Request, acc *account.Account) (error, string) {
  return sendEmail (r, acc, "activation", acc.ActivationCode)
}

func ActivationCodeForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  util.Template("activation_resend", "", map[string]interface{}{
    "uid": ps.ByName("uid"),
  })(w, r, ps)
}

func ActivationCode (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  handle := util.Handle(w, r, ps)
  if acc, err, conflict := account.GetInsecure(r.FormValue("uid")); err != nil {
    handle(err, conflict, "activation_resend", map[string]interface{}{
      "uid": r.FormValue("uid"),
    })
  } else if acc.IsActive() {
    handle(account.ErrAlreadyActivated, true, "", nil)
  } else if err, remark := activationEmail(r, acc); err != nil {
    handle(err, false, "", nil)
  } else {
    util.Template("activation_resend_success", "", map[string]interface{}{
      "name": acc.Name(),
      "uid": acc.UID,
      "remark": remark,
    })(w, r, ps)
  }
}

func passwordEmail (r *http.Request, acc *account.Account) (error, string) {
  // TODO: indicate Expires-time
  return sendEmail (r, acc, "password", acc.PasswordCode.Value)
}

func PasswordCodeForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  util.Template("passwordcode", "", map[string]interface{}{
    "uid": ps.ByName("uid"),
  })(w, r, ps)
}

func PasswordCode (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  handle := util.Handle(w, r, ps)
  if acc, err, conflict := account.GetInsecure(r.FormValue("uid")); err != nil {
    handle(err, conflict, "passwordcode", map[string]interface{}{
      "uid": r.FormValue("uid"),
    })
  } else if err := acc.CreatePasswordCode(); err != nil {
    handle(err, false, "", nil)
  } else if err, remark := passwordEmail(r, acc); err != nil {
    handle(err, false, "", nil)
  } else {
    util.Template("passwordcode_success", "", map[string]interface{}{
      "name": acc.Name(),
      "remark": remark,
    })(w, r, ps)
  }
}
