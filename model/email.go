package model

import (
	"github.com/wscherphof/entity"
	"github.com/wscherphof/essix/util"
	"strings"
)

func init() {
	entity.Register(initEmail(""))
}

type email struct {
	*entity.Base
}

func initEmail(address string) *email {
	address = strings.ToLower(address)
	address = strings.TrimSpace(address)
	return &email{Base: &entity.Base{
		ID: address,
	}}
}

/*
CreateEmailToken generates a token that is needed to change the Account's email
address.
*/
func (a *Account) CreateEmailToken(newEmail string) (err error, conflict bool) {
	var empty bool
	email := initEmail(newEmail)
	if err, empty = email.Read(email); err != nil {
		if empty {
			a.NewEmail = email.ID
			a.EmailToken = util.NewToken()
			err = a.Update(a)
		}
	} else {
		err, conflict = ErrEmailTaken, true
	}
	return
}

/*
ClearEmailToken clears the token to cancel the email address changing process.
*/
func (a *Account) ClearEmailToken(token string) (err error, conflict bool) {
	if token == "" || a.EmailToken != token {
		err, conflict = ErrInvalidCredentials, true
	} else {
		a.NewEmail = ""
		a.EmailToken = ""
		err = a.Update(a)
	}
	return
}

/*
ChangeEmail sets the Account's Email to the NewEmail, if the given token is
correct.
*/
func (a *Account) ChangeEmail(token string) (err error, conflict bool) {
	email := initEmail(a.NewEmail)
	if token == "" || a.EmailToken != token || a.NewEmail == "" {
		err, conflict = ErrInvalidCredentials, true
	} else if err, conflict = email.Create(email); err != nil {
		if conflict {
			err = ErrEmailTaken
		}
	} else {
		email.ID = a.Email
		a.Email = a.NewEmail
		a.EmailToken = ""
		if err = email.Delete(email); err == nil {
			if err = a.Update(a); err != nil {
				email.Create(email)
			}
		}
	}
	return
}
