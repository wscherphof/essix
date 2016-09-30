package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model/account"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/secure"
	"net/http"
)

func terminateEmail(r *http.Request, acc *account.Account) (err error, remark string) {
	return sendEmail(r, acc.UID, acc.Name(), "terminate", acc.TerminateCode, "")
}

func TerminateCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_ = Authentication(w, r)
	template.Run(w, r, "secure", "terminatecode", "", nil)
}

func TerminateCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	sure := r.FormValue("sure")
	if e, conflict := acc.CreateTerminateCode((sure == "affirmative")); e != nil {
		template.Error(w, r, e, conflict)
	} else if e, remark := terminateEmail(r, acc); e != nil {
		template.Error(w, r, e, false)
	} else {
		secure.Update(w, r, acc)
		template.Run(w, r, "secure", "terminatecode_success", "", map[string]interface{}{
			"Name":   acc.Name(),
			"Remark": remark,
		})
	}
}

func TerminateForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	code, cancel := r.FormValue("code"), r.FormValue("cancel")
	if cancel == "true" {
		acc.ClearTerminateCode(code)
		secure.Update(w, r, acc)
		template.Run(w, r, "secure", "terminatecode_cancelled", "", nil)
	} else {
		template.Run(w, r, "secure", "terminate", "", map[string]interface{}{
			"Account": acc,
		})
	}
}

func Terminate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	code, sure := r.FormValue("code"), r.FormValue("sure")
	if e, conflict := acc.Terminate(code, (sure == "affirmative")); e != nil {
		template.Error(w, r, e, conflict)
	} else {
		secure.LogOut(w, r, false)
		template.Run(w, r, "secure", "terminate_success", "", nil)
	}
}
