package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/secure"
	"net/http"
)

func EmailTokenForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	template.Run(w, r, "email", "EmailTokenForm", "", map[string]interface{}{
		"email": account.Email,
	})
}

func EmailToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	newEmail := r.FormValue("newemail")
	if err, conflict := account.CreateEmailToken(newEmail); err != nil {
		template.Error(w, r, err, conflict)
	} else if err, remark := sendEmail(r, account.NewEmail, "EmailToken", "/account/email?token="+account.EmailToken); err != nil {
		template.Error(w, r, err, false)
	} else {
		secure.Update(w, r, account)
		template.Run(w, r, "email", "EmailToken", "", map[string]interface{}{
			"remark": remark,
		})
	}
}

func ChangeEmailForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	token, cancel := r.FormValue("token"), r.FormValue("cancel")
	if cancel == "true" {
		if err, conflict := account.ClearEmailToken(token); err != nil {
			template.Error(w, r, err, conflict)
		} else {
			secure.Update(w, r, account)
			template.Run(w, r, "email", "ChangeEmail-cancel", "", nil)
		}
	} else {
		template.Run(w, r, "email", "ChangeEmailForm", "", map[string]interface{}{
			"account": account,
		})
	}
}

func ChangeEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	token := r.FormValue("token")
	if err, conflict := account.ChangeEmail(token); err != nil {
		template.Error(w, r, err, conflict)
	} else {
		secure.Update(w, r, account)
		template.Run(w, r, "email", "ChangeEmail", "", nil)
	}
}
