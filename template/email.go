package template

import (
	"github.com/wscherphof/essix/email"
	"github.com/wscherphof/msg"
	"net/http"
	"net/url"
)

type emailType getType

func (t *emailType) Run(recipient, subject string) (err error, message string) {
	data := make(map[string]interface{}, len(*(t.Values)) + 1)
	for key := range *t.Values {
		data[key] = t.Values.Get(key)
	}
	body := Write(t.r, t.dir, t.base, t.inner, data)
	if err = email.Send(
		msg.Translator(t.r).Get(subject),
		body,
		recipient,
	); err == email.ErrNotSentImmediately {
		err, message = nil, err.Error()
	}
	return
}

func Email(r *http.Request, dir, base string, inner ...string) *emailType {
	values, _ := url.ParseQuery("")
	return &emailType{
		baseType: &baseType{&values, nil, r},
		dir:      dir,
		base:     base,
		inner:    opt(inner...),
	}
}
