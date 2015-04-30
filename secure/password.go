package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/expeertise/model/account"
  "time"
  "errors"
)

const PWD_CODE_TIMEOUT time.Duration = 1 * time.Hour
var ErrPasswordCodeTimedOut = errors.New("Password code has timed out")

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
