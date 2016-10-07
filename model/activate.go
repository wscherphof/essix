package model

import (
	"errors"
)

var (
	ErrAlreadyActivated = errors.New("ErrAlreadyActivated")
)

func (a *Account) Activate(code string) (err error, conflict bool) {
	if a.IsActive() {
		return ErrAlreadyActivated, true
	}
	if code != a.ActivateCode {
		return ErrInvalidCredentials, true
	}
	a.ActivateCode = ""
	return a.Update(a), false
}
