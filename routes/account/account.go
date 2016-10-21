/*
Package account provides route Handles for managing accounts.

Sign up, verify email address to activate account, resend activation email,
reset password, log in, change email address verifying the new email address,
suspend account through email confirmation.

The account details are kept in the database (model.Account). On log in, a copy
of the account data is stored in an encrypted session cookie to authenticate
requests for secured resources.
*/
package account

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/secure"
	"github.com/wscherphof/essix/template"
	"net/http"
)

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
	account := secure.Authentication(r)
	t := template.GET(w, r, "account", "EditAccount")
	t.Set("email", account.Email)
	t.Run()
}
