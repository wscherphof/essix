package secure

import (
  "net/http"
  "github.com/wscherphof/secure"
  "github.com/wscherphof/msg"
  "github.com/wscherphof/expeertise/util"
  "github.com/wscherphof/expeertise/model/account"
  "github.com/wscherphof/expeertise/email"
)

func Init () {
  DefineMessages()
  secure.Init(account.Account{}, &secureDB{}, validate)
}

func sendEmail (r *http.Request, acc *account.Account, resource, code, extra string) (err error, remark string) {
  subject := msg.Msg(r)(resource + " subject")
  scheme := "http"
  if r.TLS != nil {
    scheme = "https"
  }
  path := scheme + "://" + r.Host + "/account/" + resource + "/" + acc.UID
  // TODO: format links as "buttons" instead of hyperlinks
  body := util.BTemplate(resource + "_email", "lang", map[string]interface{}{
    "action": path + "?code=" + code + "&extra=" + string(util.URLEncode([]byte(extra))),
    "cancel": path,
    "name": acc.Name(),
    "extra": extra,
  })(r)
  if e := email.Send(subject, string(body), acc.UID); e != nil {
    if e == email.ErrNotSentImmediately {
      remark = e.Error()
    } else {
      err = e
    }
  }
  return
}
