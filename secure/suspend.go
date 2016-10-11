package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/secure"
	"net/http"
)

func SuspendTokenForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_ = Authentication(w, r)
	template.Run(w, r, "suspend", "SuspendTokenForm", "", nil)
}

func SuspendToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	sure := r.FormValue("sure")
	if err, conflict := account.CreateSuspendToken((sure == "affirmative")); err != nil {
		template.Error(w, r, err, conflict)
	} else if err, remark := sendEmail(r, account.Email, "SuspendToken", "/account/suspend?token="+account.SuspendToken); err != nil {
		template.Error(w, r, err, false)
	} else {
		secure.Update(w, r, account)
		template.Run(w, r, "suspend", "SuspendToken", "", map[string]interface{}{
			"remark": remark,
		})
	}
}

func SuspendForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	token, cancel := r.FormValue("token"), r.FormValue("cancel")
	if cancel == "true" {
		account.ClearSuspendToken(token)
		secure.Update(w, r, account)
		template.Run(w, r, "suspend", "Suspend-cancel", "", nil)
	} else {
		template.Run(w, r, "suspend", "SuspendForm", "", map[string]interface{}{
			"account": account,
		})
	}
}

func Suspend(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	token, sure := r.FormValue("token"), r.FormValue("sure")
	if err, conflict := account.Suspend(token, (sure == "affirmative")); err != nil {
		template.Error(w, r, err, conflict)
	} else {
		secure.LogOut(w, r, false)
		template.Run(w, r, "suspend", "Suspend", "", nil)
	}
}
