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

func TerminateCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	_ = Authentication(w, r)
	return router.Template("secure", "terminatecode", "", nil)(w, r, ps)
}

func TerminateCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	acc := Authentication(w, r)
	sure := r.FormValue("sure")
	if e, conflict := acc.CreateTerminateCode((sure == "affirmative")); e != nil {
		err = router.NewError(e)
		err.Conflict = conflict
	} else if e, remark := terminateEmail(r, acc); e != nil {
		err = router.NewError(e)
	} else {
		secure.Update(w, r, acc)
		router.Template("secure", "terminatecode_success", "", map[string]interface{}{
			"Name":   acc.Name(),
			"Remark": remark,
		})(w, r, ps)
	}
	return
}

func TerminateForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
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
	return
}

func Terminate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	acc := Authentication(w, r)
	code, sure := r.FormValue("code"), r.FormValue("sure")
	if e, conflict := acc.Terminate(code, (sure == "affirmative")); e != nil {
		err = router.NewError(e)
		err.Conflict = conflict
	} else {
		secure.LogOut(w, r, false)
		router.Template("secure", "terminate_success", "", nil)(w, r, ps)
	}
	return
}
