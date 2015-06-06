package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/expeertise/model/account"
	"github.com/wscherphof/expeertise/router"
	"github.com/wscherphof/secure"
	"net/http"
)

func Authentication(w http.ResponseWriter, r *http.Request, optional ...bool) (ret *account.Account) {
	if auth := secure.Authentication(w, r, optional...); auth != nil {
		acc := auth.(account.Account)
		ret = &acc
	}
	return
}

func IfSecureHandle(authenticated router.ErrorHandle, unauthenticated router.ErrorHandle) router.ErrorHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) *router.Error {
		if secure.Authentication(w, r, true) != nil {
			return authenticated(w, r, ps)
		} else {
			return unauthenticated(w, r, ps)
		}
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
