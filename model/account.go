package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials   = errors.New("ErrInvalidCredentials")
	ErrPasswordEmpty        = errors.New("ErrPasswordEmpty")
	ErrPasswordsNotEqual    = errors.New("ErrPasswordsNotEqual")
	ErrEmailTaken           = errors.New("ErrEmailTaken")
	ErrNotActivated         = errors.New("ErrNotActivated")
	ErrAlreadyActivated     = errors.New("ErrAlreadyActivated")
	ErrCodeUnset            = errors.New("ErrCodeUnset")
	ErrCodeIncorrect        = errors.New("ErrCodeIncorrect")
	ErrPasswordCodeTimedOut = errors.New("ErrPasswordCodeTimedOut")
)

const (
	pwdCodeTimeOut = 1 * time.Hour
)

type password struct {
	Created time.Time
	Value   []byte
}

func newPassword(pwd1, pwd2 string) (pwd *password, err error, conflict bool) {
	if pwd1 == "" {
		err, conflict = ErrPasswordEmpty, true
	} else if pwd1 != pwd2 {
		err, conflict = ErrPasswordsNotEqual, true
	} else if hash, e := bcrypt.GenerateFromPassword([]byte(pwd1), bcrypt.DefaultCost); e != nil {
		err, conflict = e, false
	} else {
		pwd = &password{
			Created: time.Now(),
			Value:   hash,
		}
	}
	return
}

type passwordCode struct {
	Expires time.Time
	Value   string
}

type Account struct {
	*Entity
	PWD              *password
	ActivationCode   string
	PasswordCode     *passwordCode
	EmailAddressCode string
	NewUID           string
	TerminateCode    string
}

func initAccount(uid string) (account *Account) {
	account = &Account{Entity: &Entity{}}
	account.Init("account", strings.ToLower(uid))
	return
}

func NewAccount(uid, pwd1, pwd2 string) (account *Account, err error, conflict bool) {
	acc := initAccount(uid)
	acc.ActivationCode = Code()
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

func (a *Account) CreatePasswordCode() error {
	a.PasswordCode = &passwordCode{
		Expires: time.Now().Add(pwdCodeTimeOut),
		Value:   Code(),
	}
	return a.Update(a)
}

func ClearPasswordCode(uid, code string) {
	if acc, _, _ := getAccount(uid); acc != nil {
		if acc.PasswordCode.Value == code {
			acc.PasswordCode = nil
			acc.Update(acc)
		}
	}
}

func (a *Account) ChangePassword(code, pwd1, pwd2 string) (err error, conflict bool) {
	if a.PasswordCode == nil {
		err, conflict = ErrCodeUnset, true
	} else if time.Now().After(a.PasswordCode.Expires) {
		a.PasswordCode = nil
		a.Update(a)
		err, conflict = ErrPasswordCodeTimedOut, true
	} else if code == "" || code != a.PasswordCode.Value {
		err, conflict = ErrCodeIncorrect, true
	} else if pwd, e, c := newPassword(pwd1, pwd2); e != nil {
		err, conflict = e, c
	} else {
		a.PasswordCode = nil
		a.PWD = pwd
		err = a.Update(a)
	}
	return
}

func (a *Account) CreateEmailAddressCode(newUID string) error {
	a.NewUID = newUID
	a.EmailAddressCode = Code()
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

func (a *Account) CreateTerminateCode(sure bool) (err error, conflict bool) {
	if !sure {
		err, conflict = ErrCodeIncorrect, true
	} else {
		a.TerminateCode = Code()
		err = a.Update(a)
	}
	return
}

func (a *Account) ClearTerminateCode(code string) (err error) {
	if a.TerminateCode == code {
		a.TerminateCode = ""
		err = a.Update(a)
	}
	return
}

func (a *Account) Terminate(code string, sure bool) (err error, conflict bool) {
	uid := a.ID
	if !sure {
		err, conflict = ErrCodeIncorrect, true
	} else if acc, e, c := getAccount(uid); e != nil {
		err, conflict = e, c
	} else if acc.TerminateCode == "" {
		err, conflict = ErrCodeUnset, true
	} else if code == "" || code != acc.TerminateCode {
		err, conflict = ErrCodeIncorrect, true
	} else {
		err = acc.Delete()
	}
	return
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
