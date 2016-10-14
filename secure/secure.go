/*
Package secure provides authentication for the application.

The account details are kept in the database (model.Account). On log in, a copy
of the account data is stored in an encrypted session cookie to authenticate
requests for secured resources.
*/
package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/secure/db"
	"github.com/wscherphof/secure"
	"net/http"
)

/*
Authentication returns the client's account data from the encrypted cookie,
which is regularly validated with the account's record in the database.
*/
func Authentication(w http.ResponseWriter, r *http.Request, optional ...bool) (ret *model.Account) {
	if auth := secure.Authentication(w, r, optional...); auth != nil {
		acc := auth.(model.Account)
		ret = &acc
	}
	return
}

/*
secure.Handle ensures the client is logged in when accessing a certian route,
redirecting to the log in page if not. The given Handle function should call
Authentication() to get the client's account details.
*/
func Handle(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if secure.Authentication(w, r, false) != nil {
			handle(w, r, ps)
		}
	}
}

/*
secure.IfHandle calls one Hanlde function for logged-in clients, and another for
logged-out clients.
*/
func IfHandle(authenticated httprouter.Handle, unauthenticated httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if secure.Authentication(w, r, true) != nil {
			authenticated(w, r, ps)
		} else {
			unauthenticated(w, r, ps)
		}
	}
}

func init() {
	// Authentication will be based on a record of model/account
	var record = model.Account{}
	// Security keys will be found through an instance of our secureDB implementation of the secure.DB interface
	var secureDB = db.New()
	// The validate function will test whether the session still valid
	var validate = func(src interface{}) (dst interface{}, valid bool) {
		switch account := src.(type) {
		case model.Account:
			valid = account.Refresh()
			dst = account
		}
		return
	}
	secure.Configure(record, secureDB, validate)
}
