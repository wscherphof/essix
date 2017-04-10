/*
Package model manages the business data model entities and their behaviuor.
*/
package model

import (
	"errors"
	"github.com/wscherphof/entity"
	"github.com/wscherphof/essix/util"
)

var (
	ErrInvalidCredentials = errors.New("ErrInvalidCredentials")
	ErrEmailTaken         = errors.New("ErrEmailTaken")
)

/*
Account represents the user's tollgate to the application.

When logged in, an encrypted representation of the user's Account is kept in an
HTTP cookie to authorise any subsequent requests.
*/
type Account struct {
	*entity.Base
	Email         string
	Password      *password
	ActivateToken string
	PasswordToken *passwordToken
	EmailToken    string
	NewEmail      string
	SuspendToken  string
}

func init() {
	entity.Register(&Account{}).Index("Email")
}

func initAccount(id ...string) (account *Account) {
	account = &Account{Base: &entity.Base{}}
	if len(id) == 1 {
		account.ID = id[0]
	}
	return
}

/*
NewAccount creates a new Account in the database, if a unique email address is
given, and the same password twice. The password is never saved, so it's not
deductable from the database. Only a cryptographic hash is stored, to compare to
the computed hash of the password to validate on log in.
*/
func NewAccount(address, pwd1, pwd2 string) (account *Account, err error, conflict bool) {
	account = initAccount()
	account.ActivateToken = util.NewToken()
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

/*
IsActive returns wehther the account is activated to confirm the user's control
over the email address provided.
*/
func (a *Account) IsActive() bool {
	return a.ActivateToken == ""
}

/*
Refresh updates the account's field values & returns the validity of the session.
It's called by the function that validates the session cookie.
*/
func (a *Account) Refresh() (current bool) {
	if saved, err, _ := GetAccount(a.ID); err == nil {
		*a = *saved
		current = a.Password.Created.Equal(saved.Password.Created)
	}
	return
}

/*
GetAccount returns the stored Account with the given id and/or email address.
*/
func GetAccount(id string, address ...string) (account *Account, err error, conflict bool) {
	var empty bool
	account = initAccount(id)
	if len(address) == 1 && id == "" {
		index := account.Index(account, "Email")
		email := initEmail(address[0])
		err, empty = index.Read(email.ID, account)
	} else {
		err, empty = account.Read(account)
	}
	if empty {
		err, conflict = ErrInvalidCredentials, true
	}
	return
}
