package template

import (
	"github.com/wscherphof/essix/email"
	"github.com/wscherphof/msg"
	"net/http"
)

type EmailType GetType

func (t *EmailType) Run(recipient, subject string) (err error, message string) {
	body := write(t.r, t.dir, t.base, t.inner, t.data)
	if err = email.Send(
		msg.Translator(t.r).Get(subject),
		body,
		recipient,
	); err == email.ErrNotSentImmediately {
		err, message = nil, err.Error()
	}
	return
}

func Email(r *http.Request, dir, base string, inner ...string) *EmailType {
	return &EmailType{
		baseType: newBaseType(nil, r),
		dir:      dir,
		base:     base,
		inner:    opt(inner...),
	}
}
