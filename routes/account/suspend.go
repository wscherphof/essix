package account

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/secure"
	"github.com/wscherphof/essix/template"
	cookie "github.com/wscherphof/secure"
	"net/http"
)

func SuspendToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := secure.Authentication(r)
	if t := template.PRG(w, r, "suspend", "SuspendToken"); t == nil {
		return
	} else if err, conflict := account.CreateSuspendToken(
		r.FormValue("sure"),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else {
		email := template.Email(r, "suspend", "SuspendToken-email", "lang")
		email.Set("link", "https://"+r.Host+"/account/suspend?token="+account.SuspendToken)
		email.Run(account.Email, "Suspend account")
		cookie.Update(w, r, account)
		t.Run()
	}
}

func SuspendForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := secure.Authentication(r)
	t := template.GET(w, r, "suspend", "SuspendForm")
	token, cancel := r.FormValue("token"), r.FormValue("cancel")
	if cancel == "true" {
		account.ClearSuspendToken(token)
		cookie.Update(w, r, account)
		template.GET(w, r, "suspend", "Suspend-cancel", "").Run()
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
		account := secure.Authentication(r)
		if err, conflict := account.Suspend(
			r.FormValue("token"),
			r.FormValue("sure"),
		); err != nil {
			template.Error(w, r, err, conflict)
		} else {
			cookie.LogOut(w, r, false)
			t.Run()
		}
	}
}
