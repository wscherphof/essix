package account

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/secure"
	"net/http"
)

// ActivateForm renders a form to enter the activate token.
func ActivateForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := template.GET(w, r, "activate", "ActivateForm")
	t.Set("id", r.FormValue("id"))
	t.Set("token", r.FormValue("token"))
	t.Set("resend_formtoken", secure.NewFormToken(r, "/account/activate/token"))
	t.Run()
}

// Activate activates the account with the given token.
func Activate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if t := template.PRG(w, r, "activate", "Activate"); t == nil {
		return
	} else if account, err, conflict := model.GetAccount(r.FormValue("id")); err != nil {
		template.Error(w, r, err, conflict)
	} else if err, conflict = account.Activate(r.FormValue("token")); err != nil {
		template.Error(w, r, err, conflict)
	} else {
		t.Run()
	}
}

// ActivateToken sends the new activate token.
func ActivateToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if t := template.PRG(w, r, "activate", "ActivateToken"); t == nil {
		return
	} else if account, err, conflict := model.GetAccount(r.FormValue("id"), r.FormValue("email")); err != nil {
		template.Error(w, r, err, conflict)
	} else if account.IsActive() {
		template.Error(w, r, model.ErrAlreadyActivated, true)
	} else {
		activateEmail(r, account)
		t.Set("id", account.ID)
		t.Run()
	}
}

func activateEmail(r *http.Request, account *model.Account) {
	email := template.Email(r, "activate", "ActivateToken-email", "lang")
	email.Set("link", "https://"+r.Host+"/account/activate?token="+account.ActivateToken+"&id="+account.ID)
	email.Run(account.Email, "Activate account")
}
