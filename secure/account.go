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
  if account := Authentication(r); account != nil {
    UpdateAccountForm(w, r, ps)
  } else {
    SignUpForm(w, r, ps)
  }
}

func UpdateAccountForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  account := Authentication(r)
  util.Template("account", "", map[string]interface{}{
    "Account": account,
    "Countries": data.Countries(),
  })(w, r, ps)
}

func UpdateAccount (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  account := Authentication(r)
  account.Country   = r.FormValue("country")
  account.Postcode  = r.FormValue("postcode")
  account.FirstName = r.FormValue("firstname")
  account.LastName  = r.FormValue("lastname")
  handle := util.Handle(w, r, ps)
  if _, err := account.Save(); err != nil {
    handle(err, false, "", nil)
  } else {
    // TODO: update account details in autorisation token
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
