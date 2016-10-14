package model

import (
	"github.com/wscherphof/essix/util"
	"log"
)

/*
CreateSuspendToken generates a token that is needed to suspend the Account, if
sure is "affirmative".
*/
func (a *Account) CreateSuspendToken(sure string) (err error, conflict bool) {
	if sure != "affirmative" {
		err, conflict = ErrInvalidCredentials, true
	} else {
		a.SuspendToken = util.NewToken()
		err = a.Update(a)
	}
	return
}

/*
ClearSuspendToken clears the token to cancel the account suspending process.
*/
func (a *Account) ClearSuspendToken(token string) (err error) {
	if a.SuspendToken == token {
		a.SuspendToken = ""
		err = a.Update(a)
	}
	return
}

/*
Suspend deletes the Account if the given token is correct and sure is
"affirmative".
*/
func (a *Account) Suspend(token string, sure string) (err error, conflict bool) {
	if sure != "affirmative" || token == "" || a.SuspendToken != token {
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
