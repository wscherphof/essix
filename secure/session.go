package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/secure"
  "github.com/wscherphof/secure/middleware"
  "github.com/wscherphof/expeertise/util2"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/dchest/captcha"
)

func LogInForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  return util2.Template("login", "", map[string]interface{}{
    "CaptchaId": captcha.New(),
  })(w, r, ps)
}

func LogIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
    err = util2.NewError(captcha.ErrNotFound, "login")
    err.Conflict = true
  } else if acc, e, conflict := account.Get(r.FormValue("uid"), r.FormValue("pwd")); err != nil {
    err = util2.NewError(e, "login")
    err.Conflict = conflict
  } else if e := secure.LogIn(w, r, acc, acc.Complete()); err != nil {
    err = util2.NewError(e, "login")
  } else if !acc.Complete() {
    middleware.SetAuthentication(r, *acc)
    err = UpdateAccountForm(w, r, ps)
  }
  return
}

func LogOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  secure.LogOut(w, r, true)
  return
}
