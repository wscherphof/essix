package secure

import (
	"github.com/wscherphof/essix/email"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/essix/util"
	"github.com/wscherphof/msg"
	"net/http"
)

func sendEmail(r *http.Request, recipient, templateName, path string, extra ...string) (err error, remark string) {
	// TODO: format links as "buttons" instead of hyperlinks
	data := map[string]interface{}{
		"link": "https://" + r.Host + path,
	}
	if len(extra) == 1 {
		data["extra"] = util.URLEncodeString(extra[0])
	}
	body := template.Write(r, "secure", templateName+"-email", "lang", data)
	subject := msg.Msg(r)(templateName + " subject")
	if err = email.Send(subject, body, recipient); err == email.ErrNotSentImmediately {
		remark = err.Error()
	}
	return
}
