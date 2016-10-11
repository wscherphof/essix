package model

import (
	"github.com/wscherphof/essix/util"
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
		err = a.Delete(a)
	}
	return
}
