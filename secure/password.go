package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/router"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/msg"
  "github.com/dchest/captcha"
  "time"
  "errors"
)

const PWD_CODE_TIMEOUT = 1 * time.Hour
var ErrPasswordCodeTimedOut = errors.New("Password code has timed out")

func passwordEmail(r *http.Request, acc *account.Account) (error, string) {
  format := msg.Msg(r)("Time format")
  return sendEmail(r, acc, "password", acc.PasswordCode.Value, acc.PasswordCode.Expires.Format(format))
}

func PasswordCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
  return router.Template("passwordcode", "", map[string]interface{}{
    "UID": ps.ByName("uid"),
    "CaptchaId": captcha.New(),
  })(w, r, ps)
}

func PasswordCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
  uid := r.FormValue("uid")
  if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
    err = router.NewError(captcha.ErrNotFound, "passwordcode")
    err.Conflict = true
  } else if acc, e, conflict := account.GetInsecure(uid); e != nil {
    err = router.NewError(e, "passwordcode")
    err.Conflict = conflict
    err.Data = map[string]interface{}{
      "UID": uid,
    }
  } else if ! acc.IsActive() {
    err = router.NewError(account.ErrNotActivated, "activation_resend")
    err.Conflict = true
    err.Data = map[string]interface{}{
      "UID": uid,
    }
  } else if e := acc.CreatePasswordCode(PWD_CODE_TIMEOUT); e != nil {
    err = router.NewError(e)
  } else if e, remark := passwordEmail(r, acc); e != nil {
    err = router.NewError(e)
  } else {
    router.Template("passwordcode_success", "", map[string]interface{}{
      "Name": acc.Name(),
      "Remark": remark,
    })(w, r, ps)
  }
  return
}

func PasswordForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
  uid, code, extra, cancel := ps.ByName("uid"), r.FormValue("code"), r.FormValue("extra"), r.FormValue("cancel")
  expires, _ := util.URLDecode([]byte(extra))
  if cancel == "true" {
    account.ClearPasswordCode(uid, code)
    router.Template("passwordcode_cancelled", "", nil)(w, r, ps)
  } else {
    router.Template("password", "", map[string]interface{}{
      "UID": uid,
      "Code": code,
      "Expires": string(expires),
    })(w, r, ps)
  }
  return
}

func ChangePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
  uid, code, pwd1, pwd2 := r.FormValue("uid"), r.FormValue("code"), r.FormValue("pwd1"), r.FormValue("pwd2")
  if acc, e, conflict := account.GetInsecure(uid); e != nil {
    err = router.NewError(e)
    err.Conflict = conflict
  } else if acc.PasswordCode == nil {
    err = router.NewError(account.ErrPasswordCodeUnset)
    err.Conflict = true
  } else if time.Now().After(acc.PasswordCode.Expires) {
    err = router.NewError(ErrPasswordCodeTimedOut, "passwordcode")
    err.Conflict = true
    err.Data = map[string]interface{}{
      "UID": acc.UID,
    }
  } else if e, conflict := acc.ChangePassword(code, pwd1, pwd2); err != nil {
    err = router.NewError(e)
    err.Conflict = conflict
  } else {
    secure.LogOut(w, r, false)
    router.Template("password_success", "", nil)(w, r, ps)
  }
  return
}
