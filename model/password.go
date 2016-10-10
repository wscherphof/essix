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
	ErrNotActivated       = errors.New("ErrNotActivated")
	ErrPasswordTokenTimedOut = errors.New("ErrPasswordTokenTimedOut")
)

const (
	pwdTokenTimeOut = 1 * time.Hour
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

func (a *Account) ValidatePassword(password string) (err error) {
	if !a.IsActive() {
		err = ErrNotActivated
	} else if err = bcrypt.CompareHashAndPassword(a.Password.Value, []byte(password)); err != nil {
		err = ErrInvalidCredentials
	}
	return
}

type passwordToken struct {
	Expires time.Time
	Value   string
}

func (a *Account) CreatePasswordToken() error {
	a.PasswordToken = &passwordToken{
		Expires: time.Now().Add(pwdTokenTimeOut),
		Value:   util.NewToken(),
	}
	return a.Update(a)
}

func ClearPasswordCode(uid, token string) {
	if acc, _, _ := GetAccount(uid); acc != nil {
		if acc.PasswordToken.Value == token {
			acc.PasswordToken = nil
			acc.Update(acc)
		}
	}
}

func (a *Account) ChangePassword(token, pwd1, pwd2 string) (err error, conflict bool) {
	if a.PasswordToken == nil {
		err, conflict = ErrTokenUnset, true
	} else if time.Now().After(a.PasswordToken.Expires) {
		a.PasswordToken = nil
		a.Update(a)
		err, conflict = ErrPasswordTokenTimedOut, true
	} else if token == "" || token != a.PasswordToken.Value {
		err, conflict = ErrTokenIncorrect, true
	} else if pwd, e, c := newPassword(pwd1, pwd2); e != nil {
		err, conflict = e, c
	} else {
		a.PasswordToken = nil
		a.Password = pwd
		err = a.Update(a)
	}
	return
}
