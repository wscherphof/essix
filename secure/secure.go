package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/expeertise/model/account"
	"github.com/wscherphof/expeertise/router"
	"github.com/wscherphof/secure"
	"net/http"
)

func Authentication(w http.ResponseWriter, r *http.Request) (ret *account.Account) {
	if auth := secure.Authentication(w, r); auth != nil {
		acc := auth.(account.Account)
		ret = &acc
	}
	return
}

func SecureHandle(handle router.ErrorHandle) router.ErrorHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
		if Authentication(w, r) != nil {
			err = handle(w, r, ps)
		} else {
			secure.Challenge(w, r)
		}
		return
	}
}

func IfSecureHandle(authenticated router.ErrorHandle, unauthenticated router.ErrorHandle) router.ErrorHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
		if Authentication(w, r) != nil {
			err = authenticated(w, r, ps)
		} else {
			err = unauthenticated(w, r, ps)
		}
		return
	}
}

func init() {
	secure.Configure(account.Account{}, &secureDB{}, func(src interface{}) (dst interface{}, valid bool) {
		if src != nil {
			acc := src.(account.Account)
			valid = acc.Refresh()
			dst = acc
		}
		return
	})
}
