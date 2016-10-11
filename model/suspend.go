package model

import (
	"github.com/wscherphof/essix/util"
	"log"
)

func (a *Account) CreateSuspendToken(sure bool) (err error, conflict bool) {
	if !sure {
		err, conflict = ErrInvalidCredentials, true
	} else {
		a.SuspendToken = util.NewToken()
		err = a.Update(a)
	}
	return
}

func (a *Account) ClearSuspendToken(token string) (err error) {
	if a.SuspendToken == token {
		a.SuspendToken = ""
		err = a.Update(a)
	}
	return
}

func (a *Account) Suspend(token string, sure bool) (err error, conflict bool) {
	if token == "" || a.SuspendToken != token {
		err, conflict = ErrInvalidCredentials, true
	} else {
		email := initEmail(a.Email)
		if err = email.Delete(email); err == nil {
			err = a.Delete(a)
		}
		if err != nil {
			log.Printf("ERROR: Suspend account %+v", err)
		}
	}
	return
}
