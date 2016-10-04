package email

import (
	"crypto/tls"
	"errors"
	"github.com/jordan-wright/email"
	"github.com/wscherphof/essix/certs"
	"github.com/wscherphof/essix/config"
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

var ErrNotSentImmediately = errors.New("ErrNotSentImmediately")

var (
	from      string
	endpoint  string
	auth      smtp.Auth
	tlsConfig *tls.Config
)

func init() {
	store := &emailConfigStore{
		Key: "email",
		Value: &emailConfig{
			EmailAddress: "essix@gmail.com",
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
			log.Println("WARNING: email.init() stored a sample email config in DB as a template to fill manually. Restart the server to read it in.")
			log.Println("INFO: r.db('essix').table('config').get('email').update({Value: {EmailAddress: 'essix@gmail.com', PWD: 'xxx', PortNumber: '587', SmtpServer: 'smtp.gmail.com'}})")
			log.Println("INFO: (note that in gmail, you need to turn on 'Allow Less Secure Apps to Access Account' through https://myaccount.google.com/u/1/security)")
		}
	} else {
		conf := store.Value
		from = conf.EmailAddress
		endpoint = conf.SmtpServer + ":" + conf.PortNumber
		auth = smtp.PlainAuth("", conf.EmailAddress, conf.PWD, conf.SmtpServer)
		tlsConfig = certs.NewConfig(conf.SmtpServer)
	}
	initQueue()
}

func Send(subject, message string, recipients ...string) (err error) {
	if e := send(subject, message, recipients...); e != nil {
		log.Println("WARNING: error sending email, enqueueing...", e)
		err = ErrNotSentImmediately
		if e := enQueue(subject, message, recipients...); e != nil {
			err = e
		}
	}
	return
}

func send(subject, message string, recipients ...string) error {
	mail := email.NewEmail()
	mail.From = from
	mail.To = recipients
	mail.Subject = subject
	mail.HTML = []byte(message)
	return mail.SendWithTLS(endpoint, auth, tlsConfig)
}
