package secure

import (
  "net/http"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/wscherphof/secure/httprouter/middleware"
)

var Authenticate = middleware.Authenticate

func Authentication (r *http.Request) (*account.Account) {
  auth := middleware.Authentication(r).(account.Account)
  return &auth
}
