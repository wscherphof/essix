package secure

import (
	"github.com/wscherphof/essix/email"
	"github.com/wscherphof/essix/util"
	"github.com/wscherphof/msg"
	"net/http"
)

func sendEmail(r *http.Request, address, name, resource, code, extra string) (err error, remark string) {
	subject := msg.Msg(r)(resource+" subject")
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	path := scheme + "://" + r.Host + "/account/" + resource + "/" + address
	action := path + "?code=" + code + "&extra=" + string(util.URLEncode([]byte(extra)))
	// TODO: format links as "buttons" instead of hyperlinks
	body := util.BTemplate(r, "secure", resource+"_email", "lang", map[string]interface{}{
		"action": action,
		"cancel": action + "&cancel=true",
		"name":   name,
		"extra":  extra,
	})
	if e := email.Send(subject, string(body), address); e != nil {
		if e == email.ErrNotSentImmediately {
			remark = e.Error()
		} else {
			err = e
		}
	}
	return
}
