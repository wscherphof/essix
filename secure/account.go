package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/wscherphof/expeertise/data"
  "github.com/dchest/captcha"
)

func AccountForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  if acc := Authentication(r); acc != nil {
    UpdateAccountForm(w, r, ps)
  } else {
    SignUpForm(w, r, ps)
  }
}

func UpdateAccountForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  acc := Authentication(r)
  util.Template("account", "", map[string]interface{}{
    "Account": acc,
    "Countries": data.Countries(),
  })(w, r, ps)
}

func UpdateAccount (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  acc := Authentication(r)
  acc.Country   = r.FormValue("country")
  acc.Postcode  = r.FormValue("postcode")
  acc.FirstName = r.FormValue("firstname")
  acc.LastName  = r.FormValue("lastname")
  handle := util.Handle(w, r, ps)
  if err := acc.Save(); err != nil {
    handle(err, false, "", nil)
  } else {
    UpdateAuthentication(w, r, acc)
    http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
  }
}

func SignUpForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  util.Template("signup", "", map[string]interface{}{
    "Countries": data.Countries(),
    "CaptchaId": captcha.New(),
  })(w, r, ps)
}

func SignUp (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  handle := util.Handle(w, r, ps)
  if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
    handle(captcha.ErrNotFound, true, "signup", nil)
  // TODO: &account.Account{...}.Create(pwd1, pwd2)
  } else if acc, err, conflict := account.New(r.FormValue); err != nil {
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
