package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model/account"
	"github.com/wscherphof/essix/util"
	"github.com/wscherphof/essix/router"
	"github.com/wscherphof/secure"
	"net/http"
)

func emailAddressEmail(r *http.Request, acc *account.Account) (err error, remark string) {
	return sendEmail(r, acc.NewUID, acc.Name(), "emailaddress", acc.EmailAddressCode, "")
}

func EmailAddressCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	util.Template(w, r, "secure", "emailaddresscode", "", map[string]interface{}{
		"UID": acc.UID,
	})
}

func EmailAddressCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	newUID := r.FormValue("newuid")
	if e := acc.CreateEmailAddressCode(newUID); e != nil {
		router.Error(e, false)(w, r, ps)
	} else if e, remark := emailAddressEmail(r, acc); e != nil {
		router.Error(e, false)(w, r, ps)
	} else {
		secure.Update(w, r, acc)
		util.Template(w, r, "secure", "emailaddresscode_success", "", map[string]interface{}{
			"Name":   acc.Name(),
			"Remark": remark,
		})
	}
}

func EmailAddressForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	code, cancel := r.FormValue("code"), r.FormValue("cancel")
	if cancel == "true" {
		acc.ClearEmailAddressCode(code)
		secure.Update(w, r, acc)
	} else {
		util.Template(w, r, "secure", "emailaddress", "", map[string]interface{}{
			"Account": acc,
		})
	}
}

func ChangeEmailAddress(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	code := r.FormValue("code")
	if e, conflict := acc.ChangeEmailAddress(code); e != nil {
		router.Error(e, conflict)(w, r, ps)
	} else {
		secure.Update(w, r, acc)
		util.Template(w, r, "secure", "emailaddress_success", "", nil)
	}
}
