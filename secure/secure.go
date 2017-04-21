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
	db := &secureDB{Base: &entity.Base{
		ID: "secure",
		Table: "config",
	}}
	entity.Register(db)
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

Call from a Handle wrapped in secure.Handle or secure.IfHandle.
*/
func Authentication(r *http.Request) *model.Account {
	auth := secure.Authentication(r)
	acc := auth.(model.Account)
	return &acc
}
