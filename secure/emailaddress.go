package secure

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/router"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/wscherphof/secure"
)

func emailAddressEmail(r *http.Request, acc *account.Account) (error, string) {
  acc.UID = acc.NewUID
  return sendEmail(r, acc, "emailaddress", acc.EmailAddressCode, "")
}

func EmailAddressCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
  acc := Authentication(r)
  return router.Template("emailaddresscode", "", map[string]interface{}{
    "UID": acc.UID,
  })(w, r, ps)
}

func EmailAddressCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
  acc := Authentication(r)
  newUID := r.FormValue("newuid")
  if e := acc.CreateEmailAddressCode(newUID); e != nil {
    err = router.NewError(e)
  } else if e, remark := emailAddressEmail(r, acc); e != nil {
    err = router.NewError(e)
  } else {
    secure.LogIn(w, r, acc, false)
    router.Template("emailaddresscode_success", "", map[string]interface{}{
      "Name": acc.Name(),
      "Remark": remark,
    })(w, r, ps)
  }
  return
}

func EmailAddressForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
  acc := Authentication(r)
  code, cancel := r.FormValue("code"), r.FormValue("cancel")
  if cancel == "true" {
    acc.ClearEmailAddressCode(code)
    router.Template("emailaddresscode_cancelled", "", nil)(w, r, ps)
  } else {
    router.Template("emailaddress", "", map[string]interface{}{
      "Account": acc,
    })(w, r, ps)
  }
  return
}

func ChangeEmailAddress(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
  acc := Authentication(r)
  code := r.FormValue("code")
  if acc.EmailAddressCode == "" {
    err = router.NewError(account.ErrEmailAddressCodeUnset)
    err.Conflict = true
  } else if e, conflict := acc.ChangeEmailAddress(code); err != nil {
    err = router.NewError(e)
    err.Conflict = conflict
  } else {
    secure.LogIn(w, r, acc, false)
    router.Template("emailaddress_success", "", nil)(w, r, ps)
  }
  return
}
