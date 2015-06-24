package secure

import (
	"github.com/dchest/captcha"
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/expeertise/data"
	"github.com/wscherphof/expeertise/model/account"
	"github.com/wscherphof/expeertise/router"
	"github.com/wscherphof/secure"
	"net/http"
	"strings"
)

func SignUpForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	return router.Template("secure", "signup", "", map[string]interface{}{
		"Countries": data.Countries(),
		"CaptchaId": captcha.New(),
	})(w, r, ps)
}

func SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		err = router.NewError(captcha.ErrNotFound)
		err.Conflict = true
	} else if acc, e, conflict := account.New(r.FormValue("uid"), r.FormValue("pwd1"), r.FormValue("pwd2")); e != nil {
		err = router.NewError(e)
		err.Conflict = conflict
	} else if e, remark := activationEmail(r, acc); e != nil {
		err = router.NewError(e)
	} else {
		router.Template("secure", "signup_success", "", map[string]interface{}{
			"uid":    acc.UID,
			"name":   acc.Name(),
			"remark": remark,
		})(w, r, ps)
	}
	return
}

func UpdateAccountForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	acc := Authentication(w, r)
	return router.Template("secure", "account", "", map[string]interface{}{
		"Account":   acc,
		"Countries": data.Countries(),
		"Initial":   (acc.ValidateFields() != nil),
	})(w, r, ps)
}

func UpdateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	acc := Authentication(w, r)
	initial := (acc.ValidateFields() != nil)
	acc.Country = r.FormValue("country")
	acc.Postcode = strings.ToUpper(r.FormValue("postcode"))
	acc.FirstName = r.FormValue("firstname")
	acc.LastName = r.FormValue("lastname")
	if e := acc.ValidateFields(); e != nil {
		err = router.NewError(e)
		err.Conflict = true
	} else if e := acc.Save(); e != nil {
		err = router.NewError(e)
	} else if initial {
		err = router.IfError(secure.LogIn(w, r, acc))
	} else if e := secure.Update(w, r, acc); e != nil {
		err = router.NewError(e)
	} else {
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	}
	return
}
