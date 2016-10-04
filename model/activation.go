package model

import (
	"errors"
)

var (
	ErrAlreadyActivated = errors.New("ErrAlreadyActivated")
)

func (a *Account) activate(code string) (err error) {
	if a.IsActive() {
		err = ErrAlreadyActivated
	} else if code != a.ActivationCode {
		err = ErrInvalidCredentials
	} else {
		a.ActivationCode = ""
	}
	return
}

func ActivateAccount(uid string, code string) (account *Account, err error, conflict bool) {
	if acc, e, c := getAccount(uid); e != nil {
		err, conflict = e, c
	} else if e := acc.activate(code); e != nil {
		err, conflict = e, true
	} else if e := acc.Update(acc); e != nil {
		err = e
	} else {
		account = acc
	}
	return
}
