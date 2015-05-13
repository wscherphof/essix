package secure

import (
  "github.com/wscherphof/expeertise/model/account"
  "log"
)

// TODO: replace with secure.Update(acc)

func validate (itf interface{}) valid bool {
  if itf != nil {
    src := itf.(account.Account)
    if acc, err, _ := account.GetInsecure(src.UID); err != nil {
      log.Println("WARNING: Validator() failed with:", err)
    } else {
      valid = (src.Modified == acc.Modified)
    }
  }
  return
}
