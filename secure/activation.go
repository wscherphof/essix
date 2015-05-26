package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/dchest/captcha"
)

func activationEmail (r *http.Request, acc *account.Account) (error, string) {
  return sendEmail (r, acc, "activation", acc.ActivationCode, "")
}

func ActivateForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  util.Template("activation", "", map[string]interface{}{
    "uid": ps.ByName("uid"),
    "code": r.FormValue("code"),
  })(w, r, ps)
}

func Activate (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  handle := util.Handle(w, r, ps)
  if account, err, conflict := account.Activate(r.FormValue("uid"), r.FormValue("code")); err != nil {
    handle(err, conflict, "activation", nil)
  } else {
    util.Template("activation_success", "", map[string]interface{}{
      "name": account.Name(),
    })(w, r, ps)
  }
}

func ActivationCodeForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  util.Template("activation_resend", "", map[string]interface{}{
    "uid": ps.ByName("uid"),
    "CaptchaId": captcha.New(),
  })(w, r, ps)
}

func ActivationCode (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  handle := util.Handle(w, r, ps)
  if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
    handle(captcha.ErrNotFound, true, "activation_resend", nil)
  } else if acc, err, conflict := account.GetInsecure(r.FormValue("uid")); err != nil {
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
