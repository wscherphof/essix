package secure

import (
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/model/account"
)

func Init () {
  DefineMessages()
  secure.Init(account.Account{}, &secureDB{})
}
