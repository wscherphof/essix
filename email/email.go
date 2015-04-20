package email

import (
  "net/smtp"
  "github.com/wscherphof/expeertise/config"
  "log"
)

const EMAIL_KEY string = "email"

type emailConfig struct{
  EmailAddress string
  PWD string
  SmtpServer string
  PortNumber string
}

type emailConfigStore struct{
  Key string
  Value emailConfig
}

var (
  conf *emailConfig
  auth smtp.Auth
  inited bool = false
)

func Init () {
  if inited {return}
  store := new(emailConfigStore)
  if err := config.Get(EMAIL_KEY, store); err != nil {
    log.Println("DEBUG: email.Init() Get error:", err)
    store.Key = EMAIL_KEY
    if err := config.Set(store); err != nil {
      log.Println("DEBUG: email.Init() Set error:", err)
    }
  } else {
    conf = &store.Value
    auth = smtp.PlainAuth("", conf.EmailAddress, conf.PWD, conf.SmtpServer)
    inited = true
  }
}

func Send (subject, message string, recipients ...string) (error) {
  Init()
  log.Println("DEBUG: email.Send() conf:", conf)
  msg := []byte("Subject: " + subject + "\n" + message)
  endpoint := conf.SmtpServer + ":" + conf.PortNumber
  return smtp.SendMail(endpoint, auth, conf.EmailAddress, recipients, msg)
}
