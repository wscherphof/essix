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
	if token, err := ratelimit.NewToken(r); err != nil {
		template.Error(w, r, err, false)
	} else {
		template.Run(w, r, "secure", "LogInForm", "", map[string]interface{}{
			"ratelimit": token,
		})
	}
}

func LogIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if account, err, conflict := model.Challenge(r.FormValue("email"), r.FormValue("password")); err != nil {
		template.Error(w, r, err, conflict, "secure", "LogIn")
	} else if err = secure.LogIn(w, r, account); err != nil {
		template.Error(w, r, err, false, "secure", "LogIn")
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func LogOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	secure.LogOut(w, r, true)
}
