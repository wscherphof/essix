package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/expeertise/model/account"
	"github.com/wscherphof/expeertise/ratelimit"
	"github.com/wscherphof/expeertise/router"
	"github.com/wscherphof/secure"
	"net/http"
)

func LogInForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	token, e := ratelimit.NewToken(r)
	if e != nil {
		return router.NewError(e)
	}
	return router.Template("secure", "login", "", map[string]interface{}{
		"RateLimitToken": token,
	})(w, r, ps)
}

func LogIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	if acc, e, conflict := account.Get(r.FormValue("uid"), r.FormValue("pwd")); e != nil {
		err = router.NewError(e, "secure", "login")
		err.Conflict = conflict
	} else if complete := (acc.ValidateFields() == nil); complete {
		err = router.IfError(secure.LogIn(w, r, acc), "secure", "login")
	} else if e := secure.Update(w, r, acc); e != nil {
		err = router.NewError(e, "secure", "login")
	} else {
		http.Redirect(w, r, "/account", http.StatusSeeOther)
	}
	return
}

func LogOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	secure.LogOut(w, r, true)
	return
}
