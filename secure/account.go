package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/data"
	"github.com/wscherphof/essix/model/account"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/router"
	"github.com/wscherphof/secure"
	"net/http"
	"strings"
)

func SignUpForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if token, e := ratelimit.NewToken(r); e != nil {
		router.Error(e, false)(w, r, ps)
	} else {
		router.Template("secure", "signup", "", map[string]interface{}{
			"Countries":      data.Countries(),
			"RateLimitToken": token,
		})(w, r, ps)
	}
}

func SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if acc, e, conflict := account.New(r.FormValue("uid"), r.FormValue("pwd1"), r.FormValue("pwd2")); e != nil {
		router.Error(e, conflict)(w, r, ps)
	} else if e, remark := activationEmail(r, acc); e != nil {
		router.Error(e, false)(w, r, ps)
	} else {
		router.Template("secure", "signup_success", "", map[string]interface{}{
			"uid":    acc.UID,
			"name":   acc.Name(),
			"remark": remark,
		})(w, r, ps)
	}
}

func UpdateAccountForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	router.Template("secure", "account", "", map[string]interface{}{
		"Account":   acc,
		"Countries": data.Countries(),
		"Initial":   (acc.ValidateFields() != nil),
	})(w, r, ps)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := Authentication(w, r)
	initial := (acc.ValidateFields() != nil)
	acc.Country = r.FormValue("country")
	acc.Postcode = strings.ToUpper(r.FormValue("postcode"))
	acc.FirstName = r.FormValue("firstname")
	acc.LastName = r.FormValue("lastname")
	if e := acc.ValidateFields(); e != nil {
		router.Error(e, true)(w, r, ps)
	} else if e := acc.Save(); e != nil {
		router.Error(e, false)(w, r, ps)
	} else if initial {
		if e := secure.LogIn(w, r, acc); e != nil {
			router.Error(e, false)(w, r, ps)
		}
	} else if e := secure.Update(w, r, acc); e != nil {
		router.Error(e, false)(w, r, ps)
	} else {
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	}
}
