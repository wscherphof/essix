package model

import (
	"errors"
)

var (
	ErrAlreadyActivated = errors.New("ErrAlreadyActivated")
)

func (a *Account) Activate(token string) (err error, conflict bool) {
	if a.IsActive() {
		return ErrAlreadyActivated, true
	}
	if token != a.ActivateToken {
		return ErrInvalidCredentials, true
	}
	a.ActivateToken = ""
	return a.Update(a), false
}
