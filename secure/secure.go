/*
Package secure provides authentication for the application.

The account details are kept in the database (model.Account). On log in, a copy
of the account data is stored in an encrypted session cookie to authenticate
requests for secured resources.
*/
package secure

import (
	"github.com/wscherphof/entity"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/secure"
	"net/http"
)

func init() {
	var db = &secureDB{
		Base: &entity.Base{
			ID: "secure",
		},
	}
	entity.Register(db, "config")
	secure.Configure(
		model.Account{},
		db,
		validate,
	)
}

// validate tests whether the session is still valid
func validate(src interface{}) (dst interface{}, valid bool) {
	switch account := src.(type) {
	case model.Account:
		valid = account.Refresh()
		dst = account
	}
	return
}

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
