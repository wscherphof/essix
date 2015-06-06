package email

import (
	"errors"
	"github.com/wscherphof/expeertise/config"
	"log"
	"net/smtp"
)

const key = "email"

type emailConfig struct {
	EmailAddress string
	PWD          string
	SmtpServer   string
	PortNumber   string
}

type emailConfigStore struct {
	Key   string
	Value *emailConfig
}

var ErrNotSentImmediately = errors.New("Message could not be sent immediately; it's queued to get sent shortly")

var (
	conf *emailConfig
	auth smtp.Auth
)

func init() {
	store := &emailConfigStore{
		Key: "email",
		Value: &emailConfig{
			EmailAddress: "expeertise@gmail.com",
			PWD:          "",
			SmtpServer:   "smtp.gmail.com",
			PortNumber:   "587",
		},
	}
	if err := config.Get(store.Key, store); err != nil {
		log.Println("WARNING: email.init() Get error:", err)
		if err := config.Set(store); err != nil {
			log.Println("ERROR: email.init() Set error:", err)
		} else {
			log.Println("WARNING: email.init() stored a blank email config in DB as a template to fill manually")
		}
	} else {
		conf = store.Value
		auth = smtp.PlainAuth("", conf.EmailAddress, conf.PWD, conf.SmtpServer)
	}
	initQueue()
}

func Send(subject, message string, recipients ...string) (err error) {
	if e := send(subject, message, recipients...); e != nil {
		err = ErrNotSentImmediately
		if e := enQueue(subject, message, recipients...); e != nil {
			err = e
		}
	}
	return
}

const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func send(subject, message string, recipients ...string) error {
	msg := "Subject: " + subject + "\n" + mime + message
	// log.Println("DEBUG: email.send() msg:", msg)
	endpoint := conf.SmtpServer + ":" + conf.PortNumber
	return smtp.SendMail(endpoint, auth, conf.EmailAddress, recipients, []byte(msg))
}
