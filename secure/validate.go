package secure

import (
  "github.com/wscherphof/expeertise/model/account"
)

func validate (src interface{}) (dest interface{}, valid bool) {
  if src != nil {
    acc := src.(account.Account)
    valid = acc.Refresh()
    dest = acc
  }
  return
}
