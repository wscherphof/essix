package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/secure"
	"net/http"
)

func LogInForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if token, e := ratelimit.NewToken(r); e != nil {
		template.Error(w, r, e, false)
	} else {
		template.Run(w, r, "secure", "login", "", map[string]interface{}{
			"RateLimitToken": token,
		})
	}
}

func LogIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if acc, e, conflict := model.GetAccount(r.FormValue("uid"), r.FormValue("pwd")); e != nil {
		template.Error(w, r, e, conflict, "secure", "login")
	} else if e := secure.LogIn(w, r, acc); e != nil {
		template.Error(w, r, e, false, "secure", "login")
	} else {
		http.Redirect(w, r, "/account", http.StatusSeeOther)
	}
}

func LogOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	secure.LogOut(w, r, true)
}
