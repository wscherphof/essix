package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/expeertise/data"
  "github.com/wscherphof/expeertise/model/account"
  "time"
  "errors"
)

func SignUpForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  // TODO: captcha
  util.Template("signup", "", map[string]interface{}{
    "Countries": data.Countries(),
  })(w, r, ps)
}

func SignUp (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  handle := util.Handle(w, r, ps)
  if acc, err, conflict := account.New(r.FormValue); err != nil {
    handle(err, conflict, "signup", nil)
  } else if err, remark := activationEmail(r, acc); err != nil {
    handle(err, false, "signup", nil)
  } else {
    util.Template("signup_success", "", map[string]interface{}{
      "uid": acc.UID,
      "name": acc.Name(),
      "remark": remark,
    })(w, r, ps)
  }
}

func ActivateForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  util.Template("activation", "", map[string]interface{}{
    "uid": ps.ByName("uid"),
    "code": r.URL.Query().Get("code"),
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

const PWD_CODE_TIMEOUT time.Duration = 1 * time.Hour
var ErrPasswordCodeTimedOut = errors.New("Password code has timed out")

func PasswordForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  uid, code := ps.ByName("uid"), r.FormValue("code")
  if len(code) == 0 {
    account.ClearPasswordCode(uid)
    util.Template("passwordcode_cancelled", "", nil)(w, r, ps)
  } else {
    util.Template("password", "", map[string]interface{}{
      "uid": uid,
      "code": code,
    })(w, r, ps)
  }
}

func ChangePassword (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  handle := util.Handle(w, r, ps)
  if acc, err, conflict := account.GetInsecure(r.FormValue("uid")); err != nil {
    handle(err, conflict, "", nil)
  } else if acc.PasswordCode == nil {
    handle(account.ErrPasswordCodeUnset, true, "", nil)
  } else if time.Since(acc.PasswordCode.Created) > PWD_CODE_TIMEOUT {
    handle(ErrPasswordCodeTimedOut, true, "passwordcode", map[string]interface{}{
      "uid": acc.UID,
    })
  } else if err, conflict := acc.ChangePassword(r.FormValue("code"), r.FormValue("pwd1"), r.FormValue("pwd2")); err != nil {
    handle(err, conflict, "", nil)
  } else {
    util.Template("password_success", "", nil)(w, r, ps)
  }
}
