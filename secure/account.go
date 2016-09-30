package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/data"
	"github.com/wscherphof/essix/model/account"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/secure"
	"net/http"
	"strings"
)

func SignUpForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if token, e := ratelimit.NewToken(r); e != nil {
		template.Error(w, r, e, false)
	} else {
		template.Run(w, r, "secure", "signup", "", map[string]interface{}{
			"Countries":      data.Countries(),
			"RateLimitToken": token,
		})
	}
}

func SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if acc, e, conflict := account.New(r.FormValue("uid"), r.FormValue("pwd1"), r.FormValue("pwd2")); e != nil {
		template.Error(w, r, e, conflict)
	} else if e, remark := activationEmail(r, acc); e != nil {
		template.Error(w, r, e, false)
	} else {
		template.Run(w, r, "secure", "signup_success", "", map[string]interface{}{
			"uid":    acc.UID,
			"name":   acc.Name(),
			"remark": remark,
		})
	}
}

func UpdateAccountForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	template.Run(w, r, "secure", "account", "", map[string]interface{}{
		"Account":   acc,
		"Countries": data.Countries(),
		"Initial":   (acc.ValidateFields() != nil),
	})
}

func UpdateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	initial := (acc.ValidateFields() != nil)
	acc.Country = r.FormValue("country")
	acc.Postcode = strings.ToUpper(r.FormValue("postcode"))
	acc.FirstName = r.FormValue("firstname")
	acc.LastName = r.FormValue("lastname")
	if e := acc.ValidateFields(); e != nil {
		template.Error(w, r, e, true)
	} else if e := acc.Save(); e != nil {
		template.Error(w, r, e, false)
	} else if initial {
		if e := secure.LogIn(w, r, acc); e != nil {
			template.Error(w, r, e, false)
		}
	} else if e := secure.Update(w, r, acc); e != nil {
		template.Error(w, r, e, false)
	} else {
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	}
}
