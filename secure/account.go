package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/util2"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/wscherphof/expeertise/data"
  "github.com/dchest/captcha"
  "strings"
)

func UpdateAccountForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  acc := Authentication(r)
  return util2.Template("account", "", map[string]interface{}{
    "Account": acc,
    "Countries": data.Countries(),
  })(w, r, ps)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  acc := Authentication(r)
  initial := (acc.ValidateFields() != nil)
  acc.Country   = r.FormValue("country")
  acc.Postcode  = strings.ToUpper(r.FormValue("postcode"))
  acc.FirstName = r.FormValue("firstname")
  acc.LastName  = r.FormValue("lastname")
  if e := acc.ValidateFields(); e != nil {
    err = util2.NewError(e)
    err.Conflict = true
  } else if e := acc.Save(); e != nil {
    err = util2.NewError(e)
  } else if initial {
    secure.LogIn(w, r, acc, true)
  } else {
    http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
  }
  return
}

func SignUpForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  return util2.Template("signup", "", map[string]interface{}{
    "Countries": data.Countries(),
    "CaptchaId": captcha.New(),
  })(w, r, ps)
}

func SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
    err = util2.NewError(captcha.ErrNotFound, "signup")
    err.Conflict = true
  // TODO: &account.Account{...}.Create(pwd1, pwd2)
  } else if acc, e, conflict := account.New(r.FormValue); e != nil {
    err = util2.NewError(e, "signup")
    err.Conflict = conflict
  } else if e, remark := activationEmail(r, acc); e != nil {
    err = util2.NewError(e, "signup")
  } else {
    util2.Template("signup_success", "", map[string]interface{}{
      "uid": acc.UID,
      "name": acc.Name(),
      "remark": remark,
    })(w, r, ps)
  }
  return
}
