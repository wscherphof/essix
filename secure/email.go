package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/secure"
	"net/http"
)

func EmailTokenForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	t := template.GET(w, r, "email", "EmailTokenForm")
	t.Set("email", account.Email)
	t.Run()
}

func EmailToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	if t := template.PRG(w, r, "email", "EmailToken"); t == nil {
		return
	} else if err, conflict := account.CreateEmailToken(
		r.FormValue("newemail"),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else {
		email := template.Email(r, "email", "EmailToken-email", "lang")
		email.Set("link", "https://"+r.Host+"/account/email?token="+account.EmailToken)
		if err, message := email.Run(account.NewEmail, "Change email address"); err != nil {
			template.Error(w, r, err, false)
		} else {
			secure.Update(w, r, account)
			t.Set("message", message)
			t.Run()
		}
	}
}

func ChangeEmailForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	t := template.GET(w, r, "email", "ChangeEmailForm")
	token, cancel := r.FormValue("token"), r.FormValue("cancel")
	if cancel == "true" {
		if err, conflict := account.ClearEmailToken(token); err != nil {
			template.Error(w, r, err, conflict)
		} else {
			secure.Update(w, r, account)
			template.Run(w, r, "email", "ChangeEmail-cancel", "", nil)
		}
	} else {
		t.Set("email", account.Email)
		t.Set("newemail", account.NewEmail)
		t.Set("emailtoken", account.EmailToken)
		t.Run()
	}
}

func ChangeEmail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	if t := template.PRG(w, r, "email", "ChangeEmail"); t == nil {
		return
	} else if err, conflict := account.ChangeEmail(
		r.FormValue("token"),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else {
		secure.Update(w, r, account)
		t.Run()
	}
}
