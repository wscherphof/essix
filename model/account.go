package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

var (
	ErrInvalidCredentials = errors.New("ErrInvalidCredentials")
	ErrEmailTaken         = errors.New("ErrEmailTaken")
	ErrNotActivated       = errors.New("ErrNotActivated")
	ErrCodeUnset          = errors.New("ErrCodeUnset")
	ErrCodeIncorrect      = errors.New("ErrCodeIncorrect")
)

type Account struct {
	*Entity
	PWD              *password
	ActivationCode   string
	PasswordCode     *passwordCode
	EmailAddressCode string
	NewUID           string
	TerminateCode    string
}

func init() {
	Register(&Account{})
}

func initAccount(uid string) (account *Account) {
	account = &Account{Entity: &Entity{}}
	account.Init(strings.ToLower(uid))
	return
}

func NewAccount(uid, pwd1, pwd2 string) (account *Account, err error, conflict bool) {
	acc := initAccount(uid)
	acc.ActivationCode = NewCode()
	if acc.PWD, err, conflict = newPassword(pwd1, pwd2); err == nil {
		if err, conflict = acc.Create(acc); err != nil {
			if conflict {
				err = ErrEmailTaken
			}
		} else {
			account = acc
		}
	}
	return
}

func (a *Account) Name() (name string) {
	return a.ID
}

func (a *Account) IsActive() bool {
	return len(a.ActivationCode) == 0
}

func (a *Account) Refresh() (current bool) {
	if saved, e, _ := getAccount(a.ID); e == nil {
		current = a.PWD.Created.Equal(saved.PWD.Created)
		*a = *saved
	}
	return
}

func getAccount(uid string) (account *Account, err error, conflict bool) {
	acc := initAccount(uid)
	if e, found := acc.Read(acc); e != nil {
		err = e
	} else if !found {
		err, conflict = ErrInvalidCredentials, true
	} else {
		account = acc
	}
	return
}

func GetAccount(uid, pwd string) (account *Account, err error, conflict bool) {
	if acc, e, c := getAccount(uid); e != nil {
		err, conflict = e, c
	} else if !acc.IsActive() {
		err, conflict = ErrNotActivated, true
	} else if e := bcrypt.CompareHashAndPassword(acc.PWD.Value, []byte(pwd)); e != nil {
		err, conflict = ErrInvalidCredentials, true
	} else {
		pwd = ""
		account = acc
	}
	return
}

func GetAccountInsecure(uid string) (account *Account, err error, conflict bool) {
	return getAccount(uid)
}
