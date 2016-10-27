package template

import (
	"github.com/wscherphof/essix/email"
	"github.com/wscherphof/msg"
	"net/http"
)

type EmailType struct {
	*BaseType
}

/*
Run sends the formatted email.
*/
func (t *EmailType) Run(recipient, subject string) (err error, message string) {
	body := String(t.r, t.dir, t.base, t.inner(), t.data)
	if err = email.Send(
		msg.Translator(t.r).Get(subject),
		body,
		recipient,
	); err == email.ErrNotSentImmediately {
		err, message = nil, err.Error()
	}
	return
}

/*
Email sets the template to use for an email. Call Set() on the result to add
data to the template's pipeline. Call Run() to send the email.
*/
func Email(r *http.Request, dir, base string, opt_inner ...string) *EmailType {
	return &EmailType{&BaseType{nil, r, dir, base, opt_inner, nil}}
}
