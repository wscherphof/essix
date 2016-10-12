package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/template"
	"net/http"
)

func ActivateForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := template.GET(w, r, "account", "ActivateForm")
	if token, err := ratelimit.NewToken(r, "/account/activate/token"); err != nil {
		template.Error(w, r, err, false)
	} else {
		t.Set("id", r.FormValue("id"))
		t.Set("token", r.FormValue("token"))
		t.Set("ratelimit", token)
		t.Run()
	}
}

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

func ActivateTokenForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := template.GET(w, r, "account", "ActivateTokenForm")
	if token, err := ratelimit.NewToken(r); err != nil {
		template.Error(w, r, err, false)
	} else {
		t.Set("ratelimit", token)
		t.Run()
	}
}

func ActivateToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if t := template.PRG(w, r, "activate", "ActivateToken"); t == nil {
		return
	} else if account, err, conflict := model.GetAccount(r.FormValue("id"), r.FormValue("email")); err != nil {
		template.Error(w, r, err, conflict)
	} else if account.IsActive() {
		template.Error(w, r, model.ErrAlreadyActivated, true)
	} else if err, message := activateEmail(r, account); err != nil {
		template.Error(w, r, err, false)
	} else {
		t.Set("id", account.ID)
		t.Set("message", message)
		t.Run()
	}
}

func activateEmail(r *http.Request, account *model.Account) (error, string) {
	email := template.Email(r, "activate", "ActivateToken-email", "lang")
	email.Set("link", "/account/activate?token="+account.ActivateToken+"&id="+account.ID)
	return email.Run(account.Email, "Activate account")
}
