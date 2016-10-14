package model

import (
	"errors"
)

var (
	ErrAlreadyActivated = errors.New("ErrAlreadyActivated")
)

/*
Activate activates the account, or returns an error if the account was already
activated, or the token given is invalid.
*/
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
