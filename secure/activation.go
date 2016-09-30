package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model/account"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/essix/ratelimit"
	"net/http"
)

func activationEmail(r *http.Request, acc *account.Account) (error, string) {
	return sendEmail(r, acc.UID, acc.Name(), "activation", acc.ActivationCode, "")
}

func ActivateForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	template.Run(w, r, "secure", "activation", "", map[string]interface{}{
		"UID":  ps.ByName("uid"),
		"Code": r.FormValue("code"),
	})
}

func Activate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if acc, e, conflict := account.Activate(r.FormValue("uid"), r.FormValue("code")); e != nil {
		template.Error(w, r, e, conflict, "secure", "activation")
	} else {
		template.Run(w, r, "secure", "activation_success", "", map[string]interface{}{
			"Name": acc.Name(),
		})
	}
}

func ActivationCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if token, e := ratelimit.NewToken(r); e != nil {
		template.Error(w, r, e, false)
	} else {
		template.Run(w, r, "secure", "activation_resend", "", map[string]interface{}{
			"UID":            ps.ByName("uid"),
			"RateLimitToken": token,
		})
	}
}

func ActivationCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if acc, e, conflict := account.GetInsecure(r.FormValue("uid")); e != nil {
		template.Error(w, r, e, conflict)
	} else if acc.IsActive() {
		template.Error(w, r, account.ErrAlreadyActivated, true)
	} else if e, remark := activationEmail(r, acc); e != nil {
		template.Error(w, r, e, false)
	} else {
		template.Run(w, r, "secure", "activation_resend_success", "", map[string]interface{}{
			"Name":   acc.Name(),
			"UID":    acc.UID,
			"Remark": remark,
		})
	}
}
