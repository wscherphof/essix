package model

import (
	"errors"
	"github.com/wscherphof/essix/util"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrPasswordEmpty        = errors.New("ErrPasswordEmpty")
	ErrPasswordsNotEqual    = errors.New("ErrPasswordsNotEqual")
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

func (a *Account) CreatePasswordCode() error {
	a.PasswordCode = &passwordCode{
		Expires: time.Now().Add(pwdCodeTimeOut),
		Value:   util.NewToken(),
	}
	return a.Update(a)
}

func ClearPasswordCode(uid, code string) {
	if acc, _, _ := GetAccount(uid); acc != nil {
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
		a.Password = pwd
		err = a.Update(a)
	}
	return
}
