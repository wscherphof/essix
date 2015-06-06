package secure

import (
	"github.com/dchest/captcha"
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/expeertise/model/account"
	"github.com/wscherphof/expeertise/router"
	"github.com/wscherphof/secure"
	"net/http"
)

func LogInForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	return router.Template("login", "", map[string]interface{}{
		"CaptchaId": captcha.New(),
	})(w, r, ps)
}

func LogIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		err = router.NewError(captcha.ErrNotFound, "login")
		err.Conflict = true
	} else if acc, e, conflict := account.Get(r.FormValue("uid"), r.FormValue("pwd")); e != nil {
		err = router.NewError(e, "login")
		err.Conflict = conflict
	} else if complete := (acc.ValidateFields() == nil); complete {
		err = router.IfError(secure.LogIn(w, r, acc), "login")
	} else if e := secure.Update(w, r, acc); e != nil {
		err = router.NewError(e, "login")
	} else {
		http.Redirect(w, r, "/account", http.StatusSeeOther)
	}
	return
}

func LogOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	secure.LogOut(w, r, true)
	return
}
