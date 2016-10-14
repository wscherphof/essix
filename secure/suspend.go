package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/secure"
	"net/http"
)

func SuspendTokenForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_ = Authentication(w, r)
	template.GET(w, r, "suspend", "SuspendTokenForm").Run()
}

func SuspendToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	if t := template.PRG(w, r, "suspend", "SuspendToken"); t == nil {
		return
	} else if err, conflict := account.CreateSuspendToken(
		r.FormValue("sure"),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else {
		email := template.Email(r, "suspend", "SuspendToken-email", "lang")
		email.Set("link", "https://"+r.Host+"/account/suspend?token="+account.SuspendToken)
		if err, message := email.Run(account.Email, "Suspend account"); err != nil {
			template.Error(w, r, err, false)
		} else {
			secure.Update(w, r, account)
			t.Set("message", message)
			t.Run()
		}
	}
}

func SuspendForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	t := template.GET(w, r, "suspend", "SuspendForm")
	token, cancel := r.FormValue("token"), r.FormValue("cancel")
	if cancel == "true" {
		account.ClearSuspendToken(token)
		secure.Update(w, r, account)
		template.Run(w, r, "suspend", "Suspend-cancel", "")
	} else {
		t.Set("email", account.Email)
		t.Set("suspendtoken", account.SuspendToken)
		t.Run()
	}
}

func Suspend(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if t := template.PRG(w, r, "suspend", "Suspend"); t == nil {
		return
	} else {
		account := Authentication(w, r)
		if err, conflict := account.Suspend(
			r.FormValue("token"),
			r.FormValue("sure"),
		); err != nil {
			template.Error(w, r, err, conflict)
		} else {
			secure.LogOut(w, r, false)
			t.Run()
		}
	}
}
