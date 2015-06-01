package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/util2"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/dchest/captcha"
)

func activationEmail(r *http.Request, acc *account.Account) (error, string) {
  return sendEmail (r, acc, "activation", acc.ActivationCode, "")
}

func ActivateForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  return util2.Template("activation", "", map[string]interface{}{
    "UID": ps.ByName("uid"),
    "Code": r.FormValue("code"),
  })(w, r, ps)
}

func Activate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  if acc, e, conflict := account.Activate(r.FormValue("uid"), r.FormValue("code")); e != nil {
    err = util2.NewError(e, "activation")
    err.Conflict = conflict
  } else {
    util2.Template("activation_success", "", map[string]interface{}{
      "Name": acc.Name(),
    })(w, r, ps)
  }
  return
}

func ActivationCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  return util2.Template("activation_resend", "", map[string]interface{}{
    "UID": ps.ByName("uid"),
    "CaptchaId": captcha.New(),
  })(w, r, ps)
}

func ActivationCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
    err = util2.NewError(captcha.ErrNotFound, "activation_resend")
    err.Conflict = true
  } else if acc, e, conflict := account.GetInsecure(r.FormValue("uid")); e != nil {
    err = util2.NewError(e, "activation_resend")
    err.Conflict = conflict
    err.Data = map[string]interface{}{
      "UID": r.FormValue("uid"),
    }
  } else if acc.IsActive() {
    err = util2.NewError(account.ErrAlreadyActivated)
    err.Conflict = true
  } else if e, remark := activationEmail(r, acc); e != nil {
    err = util2.NewError(e)
  } else {
    util2.Template("activation_resend_success", "", map[string]interface{}{
      "Name": acc.Name(),
      "UID": acc.UID,
      "Remark": remark,
    })(w, r, ps)
  }
  return
}
