package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/expeertise/data"
  "github.com/wscherphof/expeertise/model/account"
  // TODO: get rid of the email dependency for the ErrNotSentImmediately remark
  "github.com/wscherphof/expeertise/email"
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
  util.Template("activation", "", map[string]interface{}{
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
    w.Write(util.BTemplate("activation_error-tail", "", nil)(r))
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
  if acc, err, conflict := account.GetInsecure(r.FormValue("uid")); err != nil {
    if conflict {
      util.Error(w, r, ps, err, http.StatusConflict)
    } else {
      util.Error(w, r, ps, err)
    }
  } else if acc.PasswordCode == nil {
    util.Error(w, r, ps, account.ErrPasswordCodeUnset, http.StatusConflict)
  } else if time.Since(acc.PasswordCode.Created) > PWD_CODE_TIMEOUT {
    util.Error(w, r, ps, ErrPasswordCodeTimedOut, http.StatusConflict)
    w.Write(util.BTemplate("passwordcode_error-tail", "", map[string]interface{}{
      "uid": acc.UID,
    })(r))
  } else if err, conflict := acc.ChangePassword(r.FormValue("code"), r.FormValue("pwd1"), r.FormValue("pwd2")); err != nil {
    if conflict {
      util.Error(w, r, ps, err, http.StatusConflict)
    } else {
      util.Error(w, r, ps, err)
    }
  } else {
    util.Template("password_success", "", nil)(w, r, ps)
  }
}
