package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/secure"
	"net/http"
)

func emailAddressEmail(r *http.Request, acc *model.Account) (err error, remark string) {
	return sendEmail(r, acc.NewUID, acc.Name(), "emailaddress", acc.EmailAddressCode, "")
}

func EmailAddressCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	template.Run(w, r, "secure", "emailaddresscode", "", map[string]interface{}{
		"UID": acc.ID,
	})
}

func EmailAddressCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	newUID := r.FormValue("newuid")
	if e := acc.CreateEmailAddressCode(newUID); e != nil {
		template.Error(w, r, e, false)
	} else if e, remark := emailAddressEmail(r, acc); e != nil {
		template.Error(w, r, e, false)
	} else {
		secure.Update(w, r, acc)
		template.Run(w, r, "secure", "emailaddresscode_success", "", map[string]interface{}{
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
		template.Run(w, r, "secure", "emailaddress", "", map[string]interface{}{
			"Account": acc,
		})
	}
}

func ChangeEmailAddress(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	code := r.FormValue("code")
	if e, conflict := acc.ChangeEmailAddress(code); e != nil {
		template.Error(w, r, e, conflict)
	} else {
		secure.Update(w, r, acc)
		template.Run(w, r, "secure", "emailaddress_success", "", nil)
	}
}
