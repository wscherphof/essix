package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/util2"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/wscherphof/expeertise/data"
  "github.com/dchest/captcha"
  "strings"
)

func UpdateAccountForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  acc := Authentication(r)
  util2.Template("account", "", map[string]interface{}{
    "Account": acc,
    "Countries": data.Countries(),
  })(w, r, ps)
  return
}

func UpdateAccount (w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  acc := Authentication(r)
  acc.Country   = r.FormValue("country")
  acc.Postcode  = strings.ToUpper(r.FormValue("postcode"))
  acc.FirstName = r.FormValue("firstname")
  acc.LastName  = r.FormValue("lastname")
  if e := acc.Save(); e != nil {
    err = &util2.Error{Error: e}
  } else {
    UpdateAuthentication(w, r, acc)
    http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
  }
  return
}

// TODO: sign up w/ just email & pwd; then on first login, ask further details
func SignUpForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  util2.Template("signup", "", map[string]interface{}{
    "Countries": data.Countries(),
    "CaptchaId": captcha.New(),
  })(w, r, ps)
  return
}

func SignUp (w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *util2.Error) {
  if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
    err = &util2.Error{Error: captcha.ErrNotFound, Conflict: true, Tail: "signup"}
  // TODO: &account.Account{...}.Create(pwd1, pwd2)
  } else if acc, e, conflict := account.New(r.FormValue); e != nil {
    err = &util2.Error{Error: e, Conflict: conflict, Tail: "signup"}
  } else if e, remark := activationEmail(r, acc); e != nil {
    err = &util2.Error{Error: e, Conflict: false, Tail: "signup"}
  } else {
    util2.Template("signup_success", "", map[string]interface{}{
      "uid": acc.UID,
      "name": acc.Name(),
      "remark": remark,
    })(w, r, ps)
  }
  return
}
