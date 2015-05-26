package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/wscherphof/msg"
  "github.com/dchest/captcha"
  "time"
  "errors"
)

const PWD_CODE_TIMEOUT time.Duration = 1 * time.Hour
var ErrPasswordCodeTimedOut = errors.New("Password code has timed out")

func passwordEmail (r *http.Request, acc *account.Account) (error, string) {
  format := msg.Msg(r)("Time format")
  return sendEmail (r, acc, "password", acc.PasswordCode.Value, acc.PasswordCode.Expires.Format(format))
}

func PasswordCodeForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  util.Template("passwordcode", "", map[string]interface{}{
    "uid": ps.ByName("uid"),
    "CaptchaId": captcha.New(),
  })(w, r, ps)
}

func PasswordCode (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  uid := r.FormValue("uid")
  handle := util.Handle(w, r, ps)
  if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
    handle(captcha.ErrNotFound, true, "passwordcode", nil)
  } else if acc, err, conflict := account.GetInsecure(uid); err != nil {
    handle(err, conflict, "passwordcode", map[string]interface{}{"uid": uid})
  } else if ! acc.IsActive() {
    handle(account.ErrNotActivated, true, "activation_resend", map[string]interface{}{"uid": uid})
  } else if err := acc.CreatePasswordCode(PWD_CODE_TIMEOUT); err != nil {
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

func PasswordForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  uid, code, extra := ps.ByName("uid"), r.FormValue("code"), r.FormValue("extra")
  expires, _ := util.URLDecode([]byte(extra))
  if len(code) == 0 {
    account.ClearPasswordCode(uid)
    util.Template("passwordcode_cancelled", "", nil)(w, r, ps)
  } else {
    util.Template("password", "", map[string]interface{}{
      "uid": uid,
      "code": code,
      "expires": string(expires),
    })(w, r, ps)
  }
}

func ChangePassword (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  handle := util.Handle(w, r, ps)
  if acc, err, conflict := account.GetInsecure(r.FormValue("uid")); err != nil {
    handle(err, conflict, "", nil)
  } else if acc.PasswordCode == nil {
    handle(account.ErrPasswordCodeUnset, true, "", nil)
  } else if time.Now().After(acc.PasswordCode.Expires) {
    handle(ErrPasswordCodeTimedOut, true, "passwordcode", map[string]interface{}{"uid": acc.UID})
  } else if err, conflict := acc.ChangePassword(r.FormValue("code"), r.FormValue("pwd1"), r.FormValue("pwd2")); err != nil {
    handle(err, conflict, "", nil)
  } else {
    secure.LogOut(w, r, false)
    util.Template("password_success", "", nil)(w, r, ps)
  }
}
