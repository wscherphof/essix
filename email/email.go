package email

import (
  "net/smtp"
  "github.com/wscherphof/expeertise/config"
  "log"
  "errors"
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
  Value *emailConfig
}

var (
  ErrNotSentImmideately = errors.New("Message could not be sent immediately; it's queued to get sent shortly")
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
    } else {
      log.Println("DEBUG: email.Init() stored a blank email config in DB as a template to fill manually")
    }
  } else {
    conf = store.Value
    auth = smtp.PlainAuth("", conf.EmailAddress, conf.PWD, conf.SmtpServer)
    inited = true
  }
  initQueue()
}

func Send (subject, message string, recipients ...string) (err error) {
  Init()
  if e := send(subject, message, recipients...); e != nil {
    err = ErrNotSentImmideately
    if e := enQueue(subject, message, recipients...); e != nil {
      err = e
    }
  }
  return
}

func send (subject, message string, recipients ...string) (error) {
  msg := []byte("Subject: " + subject + "\n" + message)
  endpoint := conf.SmtpServer + ":" + conf.PortNumber
  return smtp.SendMail(endpoint, auth, conf.EmailAddress, recipients, msg)
}
