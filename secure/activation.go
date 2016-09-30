package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model/account"
	"github.com/wscherphof/essix/util"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/router"
	"net/http"
)

func activationEmail(r *http.Request, acc *account.Account) (error, string) {
	return sendEmail(r, acc.UID, acc.Name(), "activation", acc.ActivationCode, "")
}

func ActivateForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	util.Template(w, r, "secure", "activation", "", map[string]interface{}{
		"UID":  ps.ByName("uid"),
		"Code": r.FormValue("code"),
	})
}

func Activate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if acc, e, conflict := account.Activate(r.FormValue("uid"), r.FormValue("code")); e != nil {
		router.Error(e, conflict, "secure", "activation")(w, r, ps)
	} else {
		util.Template(w, r, "secure", "activation_success", "", map[string]interface{}{
			"Name": acc.Name(),
		})
	}
}

func ActivationCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if token, e := ratelimit.NewToken(r); e != nil {
		router.Error(e, false)(w, r, ps)
	} else {
		util.Template(w, r, "secure", "activation_resend", "", map[string]interface{}{
			"UID":            ps.ByName("uid"),
			"RateLimitToken": token,
		})
	}
}

func ActivationCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if acc, e, conflict := account.GetInsecure(r.FormValue("uid")); e != nil {
		router.Error(e, conflict)(w, r, ps)
	} else if acc.IsActive() {
		router.Error(account.ErrAlreadyActivated, true)(w, r, ps)
	} else if e, remark := activationEmail(r, acc); e != nil {
		router.Error(e, false)(w, r, ps)
	} else {
		util.Template(w, r, "secure", "activation_resend_success", "", map[string]interface{}{
			"Name":   acc.Name(),
			"UID":    acc.UID,
			"Remark": remark,
		})
	}
}
