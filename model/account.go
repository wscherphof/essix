package model

import (
	"errors"
	"github.com/wscherphof/entity"
	"github.com/wscherphof/essix/util"
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

type Email struct {
	*entity.Base
}

type Account struct {
	*entity.Base
	Email            string
	Password         *password
	ActivateCode     string
	PasswordCode     *passwordCode
	EmailAddressCode string
	NewEmail         string
	TerminateCode    string
}

func init() {
	entity.Register(&Email{})
	entity.Register(&Account{}).Index("Email")
}

func initEmail(address string) *Email {
	return &Email{Base: &entity.Base{
		ID: strings.ToLower(address),
	}}
}

func initAccount(id ...string) (account *Account) {
	account = &Account{Base: &entity.Base{}}
	if len(id) == 1 {
		account.ID = id[0]
	}
	return
}

func NewAccount(address, pwd1, pwd2 string) (account *Account, err error, conflict bool) {
	account = initAccount()
	account.ActivateCode = util.NewToken()
	if account.Password, err, conflict = newPassword(pwd1, pwd2); err != nil {
		return
	}
	email := initEmail(address)
	if err, conflict = email.Create(email); err != nil {
		if conflict {
			err = ErrEmailTaken
		}
	} else {
		account.Email = email.ID
		err = account.Update(account)
	}
	return
}

func (a *Account) Name() (name string) {
	return a.Email
}

func (a *Account) IsActive() bool {
	return a.ActivateCode == ""
}

// Refresh updates the account's field values & returns the validity of the session
func (a *Account) Refresh() (current bool) {
	if saved, err, _ := GetAccount(a.ID); err == nil {
		*a = *saved
		current = a.Password.Created.Equal(saved.Password.Created)
	}
	return
}

func GetAccount(id string, email ...string) (account *Account, err error, conflict bool) {
	account = initAccount(id)
	if len(email) == 1 && id == "" {
		index := entity.Index(account, "Email")
		err = index.Read(email[0], account)
	} else {
		err = account.Read(account)
	}
	if err == entity.ErrEmptyResult {
		err, conflict = ErrInvalidCredentials, true
	}
	return
}

func Challenge(email string, password string) (account *Account, err error, conflict bool) {
	if account, err, conflict = GetAccount("", email); err != nil {
		return
	}
	if !account.IsActive() {
		err, conflict = ErrNotActivated, true
	} else if err = bcrypt.CompareHashAndPassword(account.Password.Value, []byte(password)); err != nil {
		err, conflict = ErrInvalidCredentials, true
	}
	return
}
