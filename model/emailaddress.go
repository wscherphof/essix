package model

import (
	"github.com/wscherphof/essix/entity"
)

func (a *Account) CreateEmailAddressCode(newUID string) error {
	a.NewUID = newUID
	a.EmailAddressCode = entity.Token()
	return a.Update(a)
}

func (a *Account) ClearEmailAddressCode(code string) (err error) {
	if a.EmailAddressCode == code {
		a.NewUID = ""
		a.EmailAddressCode = ""
		err = a.Update(a)
	}
	return
}

func (a *Account) ChangeEmailAddress(code string) (err error, conflict bool) {
	uid := a.ID
	if acc, e, c := getAccount(uid); e != nil {
		err, conflict = e, c
	} else if acc.EmailAddressCode == "" {
		err, conflict = ErrCodeUnset, true
	} else if code == "" || code != acc.EmailAddressCode {
		err, conflict = ErrCodeIncorrect, true
	} else {
		acc.ID = acc.NewUID
		if e := acc.ClearEmailAddressCode(code); e != nil {
			err, conflict = e, false
		} else {
			*a = *acc
		}
	}
	return
}
