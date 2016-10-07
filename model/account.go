package model

import (
	"errors"
	"github.com/wscherphof/entity"
	"github.com/wscherphof/essix/util"
	// "golang.org/x/crypto/bcrypt"
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
	ActivationCode   string
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
	account.ActivationCode = util.NewToken()
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
	return len(a.ActivationCode) == 0
}

// Refresh updates the account's field values & returns the validity of the session
func (a *Account) Refresh() (current bool) {
	if saved, err, _ := getAccount(a.ID); err == nil {
		*a = *saved
		current = a.Password.Created.Equal(saved.Password.Created)
	}
	return
}

func getAccount(id string) (account *Account, err error, conflict bool) {
	account = initAccount(id)
	if err = account.Read(account); err != nil {
		if err == entity.ErrEmptyResult {
			err, conflict = ErrInvalidCredentials, true
		}
	}
	return
}

// func GetAccount(uid, pwd string) (account *Account, err error, conflict bool) {
// 	if acc, e, c := getAccount(uid); e != nil {
// 		err, conflict = e, c
// 	} else if !acc.IsActive() {
// 		err, conflict = ErrNotActivated, true
// 	} else if e := bcrypt.CompareHashAndPassword(acc.PWD.Value, []byte(pwd)); e != nil {
// 		err, conflict = ErrInvalidCredentials, true
// 	} else {
// 		pwd = ""
// 		account = acc
// 	}
// 	return
// }

// func GetAccountInsecure(uid string) (account *Account, err error, conflict bool) {
// 	return getAccount(uid)
// }
