package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model/account"
	"github.com/wscherphof/essix/router"
	"github.com/wscherphof/secure"
	"net/http"
)

func terminateEmail(r *http.Request, acc *account.Account) (err error, remark string) {
	return sendEmail(r, acc.UID, acc.Name(), "terminate", acc.TerminateCode, "")
}

func TerminateCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_ = Authentication(w, r)
	router.Template("secure", "terminatecode", "", nil)(w, r, ps)
}

func TerminateCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	sure := r.FormValue("sure")
	if e, conflict := acc.CreateTerminateCode((sure == "affirmative")); e != nil {
		router.Error(e, conflict)(w, r, ps)
	} else if e, remark := terminateEmail(r, acc); e != nil {
		router.Error(e, false)(w, r, ps)
	} else {
		secure.Update(w, r, acc)
		router.Template("secure", "terminatecode_success", "", map[string]interface{}{
			"Name":   acc.Name(),
			"Remark": remark,
		})(w, r, ps)
	}
}

func TerminateForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	code, cancel := r.FormValue("code"), r.FormValue("cancel")
	if cancel == "true" {
		acc.ClearTerminateCode(code)
		secure.Update(w, r, acc)
		router.Template("secure", "terminatecode_cancelled", "", nil)(w, r, ps)
	} else {
		router.Template("secure", "terminate", "", map[string]interface{}{
			"Account": acc,
		})(w, r, ps)
	}
}

func Terminate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	code, sure := r.FormValue("code"), r.FormValue("sure")
	if e, conflict := acc.Terminate(code, (sure == "affirmative")); e != nil {
		router.Error(e, conflict)(w, r, ps)
	} else {
		secure.LogOut(w, r, false)
		router.Template("secure", "terminate_success", "", nil)(w, r, ps)
	}
}
