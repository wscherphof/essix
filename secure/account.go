package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/template"
	"net/http"
)

/*
NewAccountForm renders the form to create a new account.
*/
func NewAccountForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := template.GET(w, r, "account", "NewAccountForm")
	if token, err := ratelimit.NewToken(r); err != nil {
		template.Error(w, r, err, false)
	} else {
		t.Set("ratelimit", token)
		t.Run()
	}
}

/*
NewAccount creates a new account in the database.
*/
func NewAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if t := template.PRG(w, r, "account", "NewAccount"); t == nil {
		return
	} else if account, err, conflict := model.NewAccount(
		r.FormValue("email"),
		r.FormValue("pwd1"),
		r.FormValue("pwd2"),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else if err, message := activateEmail(r, account); err != nil {
		template.Error(w, r, err, false)
	} else {
		t.Set("id", account.ID)
		t.Set("message", message)
		t.Run()
	}
}

/*
EditAccount renders a page with links and buttons to edit a logged-in client's
account details.
*/
func EditAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	t := template.GET(w, r, "account", "EditAccount")
	t.Set("email", account.Email)
	t.Run()
}
