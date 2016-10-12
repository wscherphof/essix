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
		ensure(r),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else if err, remark := sendEmail(r, account.Email, "SuspendToken", "/account/suspend?token="+account.SuspendToken); err != nil {
		template.Error(w, r, err, false)
	} else {
		secure.Update(w, r, account)
		t.Set("remark", remark)
		t.Run()
	}
}

func SuspendForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	t := template.GET(w, r, "suspend", "SuspendForm")
	token, cancel := r.FormValue("token"), r.FormValue("cancel")
	if cancel == "true" {
		account.ClearSuspendToken(token)
		secure.Update(w, r, account)
		template.Run(w, r, "suspend", "Suspend-cancel", "", nil)
	} else {
		t.Set("email", account.Email)
		t.Set("suspendtoken", account.SuspendToken)
		t.Run()
	}
}

func Suspend(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	if t := template.PRG(w, r, "suspend", "Suspend"); t == nil {
		return
	} else if err, conflict := account.Suspend(
		r.FormValue("token"),
		ensure(r),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else {
		secure.LogOut(w, r, false)
		t.Run()
	}
}

func ensure(r *http.Request) bool {
	return r.FormValue("sure") == "affirmative"
}
