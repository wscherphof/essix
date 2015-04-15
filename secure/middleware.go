package secure

import (
  "net/http"
  "github.com/wscherphof/expeertise/model"
  "github.com/wscherphof/secure/httprouter/middleware"
)

var Authenticate = middleware.Authenticate

func Authentication (r *http.Request) (*model.Account) {
  auth := middleware.Authentication(r).(model.Account)
  return &auth
}
