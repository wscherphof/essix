package account

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/secure"
	"github.com/wscherphof/essix/template"
	cookie "github.com/wscherphof/secure"
	"net/http"
)

// EmailTokenForm renders a form to request a token to change the account's email address.
func EmailTokenForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := secure.Authentication(r)
	t := template.GET(w, r, "email", "EmailTokenForm")
	t.Set("email", account.Email)
	t.Run()
}

// EmailToken sends the token to change the account's email address.
func EmailToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := secure.Authentication(r)
	if t := template.PRG(w, r, "email", "EmailToken"); t == nil {
		return
	} else if err, conflict := account.CreateEmailToken(
		r.FormValue("newemail"),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else {
		email := template.Email(r, "email", "EmailToken-email", "lang")
		email.Set("link", "https://"+r.Host+"/account/email?token="+account.EmailToken)
		email.Run(account.NewEmail, "Change email")
		cookie.Update(w, r, account)
		t.Run()
	}
}

// ChangeEmailForm accepts the token sent to the new email address to set for the account.
func ChangeEmailForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := secure.Authentication(r)
	t := template.GET(w, r, "email", "ChangeEmailForm")
	token, cancel := r.FormValue("token"), r.FormValue("cancel")
	if cancel == "true" {
		if err, conflict := account.ClearEmailToken(token); err != nil {
			template.Error(w, r, err, conflict)
		} else {
			cookie.Update(w, r, account)
			template.GET(w, r, "email", "ChangeEmail-cancel", "").Run()
		}
	} else {
		t.Set("email", account.Email)
		t.Set("newemail", account.NewEmail)
		t.Set("emailtoken", account.EmailToken)
		t.Run()
	}
}

// ChangeEmail sets the new email address for the acoount if the given token is correct.
func ChangeEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := secure.Authentication(r)
	if t := template.PRG(w, r, "email", "ChangeEmail"); t == nil {
		return
	} else if err, conflict := account.ChangeEmail(
		r.FormValue("token"),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else {
		cookie.Update(w, r, account)
		t.Run()
	}
}
