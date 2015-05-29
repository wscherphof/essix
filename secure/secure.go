package secure

import (
  "net/http"
  "github.com/wscherphof/secure"
  "github.com/wscherphof/secure/httprouter/middleware"
  "github.com/wscherphof/expeertise/model/account"
)


var (
  Authenticated         = middleware.Authenticated
  AuthenticationHandler = middleware.AuthenticationHandler
  UpdateAuthentication  = secure.UpdateAuthentication
)

func Authentication (r *http.Request) (ret *account.Account) {
  if auth := middleware.Authentication(r); auth != nil {
    acc := auth.(account.Account)
    ret = &acc
  }
  return
}

func Init () {
  DefineMessages()
  secure.Init(account.Account{}, &secureDB{}, func (src interface{}) (dest interface{}, valid bool) {
    if src != nil {
      acc := src.(account.Account)
      valid = acc.Refresh()
      dest = acc
    }
    return
  })
}
