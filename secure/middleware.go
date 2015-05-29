package secure

import (
  "net/http"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/wscherphof/secure/httprouter/middleware"
)

var UpdateAuthentication = middleware.UpdateAuthentication

var Authenticate = middleware.Authenticate

var IfAuthenticate = middleware.IfAuthenticate

func Authentication (r *http.Request) (ret *account.Account) {
  if auth := middleware.Authentication(r); auth != nil {
    acc := auth.(account.Account)
    ret = &acc
  }
  return
}
