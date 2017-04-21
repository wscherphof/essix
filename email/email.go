/*
Package email manages outgoing email traffic.

It relies on a configuration item in the database, holding the credentials of
the sending email account, and the smtp server address and port number.

If the configuration item is missing, a template for it is created in the
database. Update it with your account's details, then restart the server to load
them.
*/
package email

import (
	"crypto/tls"
	"errors"
	"github.com/jordan-wright/email"
	"github.com/wscherphof/certs"
	"github.com/wscherphof/entity"
	"log"
	"net/smtp"
)

type config struct {
	*entity.Base
	EmailAddress string
	PWD          string
	SmtpServer   string
	PortNumber   string
}

var (
	from      string
	endpoint  string
	auth      smtp.Auth
	tlsConfig *tls.Config
	conf      = &config{
		Base: &entity.Base{
			ID: "email",
			Table: "config",
		},
		EmailAddress: "essix@gmail.com",
		PWD:          "",
		SmtpServer:   "smtp.gmail.com",
		PortNumber:   "587",
	}
	ErrNotSentImmediately = errors.New("ErrNotSentImmediately")
)

func init() {
	entity.Register(conf)
	if err, _ := conf.Read(conf); err != nil {
		log.Println("WARNING: email.init() error reading config:", err)
		if err := conf.Update(conf); err != nil {
			log.Println("ERROR: email.init() Set error:", err)
		} else {
			log.Println("WARNING: email.init() stored a sample email config in DB as a template to fill manually.")
			log.Println("INFO: After updating the email config in the database, restart the server to read it in.")
			log.Println("INFO: r.db('essix').table('config').get('email').update({EmailAddress: 'essix@gmail.com', PWD: 'xxx', PortNumber: '587', SmtpServer: 'smtp.gmail.com'})")
			log.Println("INFO: (note that in gmail, you need to turn on 'Allow Less Secure Apps to Access Account' through https://myaccount.google.com/u/1/security)")
		}
	} else {
		from = conf.EmailAddress
		endpoint = conf.SmtpServer + ":" + conf.PortNumber
		auth = smtp.PlainAuth("", conf.EmailAddress, conf.PWD, conf.SmtpServer)
		tlsConfig = certs.NewConfig(conf.SmtpServer)
	}
}

/*
SendSync sends an eamil message, or enqueues it if it couldn't be sent at once.
*/
func SendSync(subject, message string, recipients ...string) (err error) {
	if e := send(subject, message, recipients...); e != nil {
		err = ErrNotSentImmediately
		if e := enQueue(subject, message, recipients...); e != nil {
			err = e
		}
	}
	return
}

/*
Send sends an eamil message asynchronously, or enqueues it if it couldn't be
sent at once.
*/
func Send(subject, message string, recipients ...string) {
	go func() {
		if err := send(subject, message, recipients...); err != nil {
			if err := enQueue(subject, message, recipients...); err != nil {
				log.Println("ERROR: enqueueing email failed", err)
			}
		}
	}()
}

func send(subject, message string, recipients ...string) (err error) {
	mail := email.NewEmail()
	mail.From = from
	mail.To = recipients
	mail.Subject = subject
	mail.HTML = []byte(message)
	if err := mail.SendWithTLS(endpoint, auth, tlsConfig); err != nil {
		log.Println("WARNING: sending email failed", err)
	}
	return
}
