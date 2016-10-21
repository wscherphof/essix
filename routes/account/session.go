package account

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/template"
	cookie "github.com/wscherphof/secure"
	"net/http"
)

func logInError(w http.ResponseWriter, r *http.Request, err error, conflict bool, id string) {
	template.ErrorTail(w, r, err, conflict, "session", "LogIn-error-tail", map[string]interface{}{
		"id": id,
	})
}

func LogIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if account, err, conflict := model.GetAccount("", r.FormValue("email")); err != nil {
		logInError(w, r, err, conflict, account.ID)
	} else if err = account.ValidatePassword(r.FormValue("password")); err != nil {
		logInError(w, r, err, true, account.ID)
	} else if err = cookie.LogIn(w, r, account); err != nil {
		logInError(w, r, err, false, account.ID)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func LogOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookie.LogOut(w, r, true)
}
