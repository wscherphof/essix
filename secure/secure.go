package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/secure/db"
	"github.com/wscherphof/secure"
	"net/http"
)

func Authentication(w http.ResponseWriter, r *http.Request, optional ...bool) (ret *model.Account) {
	if auth := secure.Authentication(w, r, optional...); auth != nil {
		acc := auth.(model.Account)
		ret = &acc
	}
	return
}

func SecureHandle(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if secure.Authentication(w, r, false) != nil {
			handle(w, r, ps)
		}
	}
}

func IfSecureHandle(authenticated httprouter.Handle, unauthenticated httprouter.Handle) httprouter.Handle {
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
		if src != nil {
			acc := src.(model.Account)
			// Refresh updates the account's field values & returns the validity of the session
			valid = acc.Refresh()
			dst = acc
		}
		return
	}
	secure.Configure(record, secureDB, validate)
}
