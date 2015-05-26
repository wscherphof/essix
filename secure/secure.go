package secure

import (
  "github.com/wscherphof/secure"
  "github.com/wscherphof/expeertise/model/account"
)

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
