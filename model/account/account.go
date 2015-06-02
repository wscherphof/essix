package account

import (
	"errors"
	"github.com/wscherphof/expeertise/db"
	"github.com/wscherphof/expeertise/util"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials        = errors.New("Unknown email address or incorrect password or activation code")
	ErrPasswordEmpty             = errors.New("Password empty")
	ErrPasswordsNotEqual         = errors.New("Passwords not equal")
	ErrEmailTaken                = errors.New("Email address taken")
	ErrNotActivated              = errors.New("Account hasn't been activated yet")
	ErrAlreadyActivated          = errors.New("Account is already activated")
	ErrPasswordCodeUnset         = errors.New("PasswordCode is nil")
	ErrPasswordCodeIncorrect     = errors.New("Password code given is incorrect")
	ErrValidationFailed          = errors.New("Field values are missing or incorrect")
	ErrEmailAddressCodeUnset     = errors.New("Email address code is empty")
	ErrEmailAddressCodeIncorrect = errors.New("Email address code given is incorrect")
)

const ACCOUNT_TABLE = "account"

func init() {
	if cursor, _ := db.TableCreatePK(ACCOUNT_TABLE, "UID"); cursor != nil {
		log.Println("INFO: table created:", ACCOUNT_TABLE)
	}
}

type password struct {
	Created time.Time
	Value   []byte
}

func newPassword(pwd1, pwd2 string) (pwd *password, err error) {
	if pwd1 == "" {
		err = ErrPasswordEmpty
	} else if pwd1 != pwd2 {
		err = ErrPasswordsNotEqual
	} else if hash, e := bcrypt.GenerateFromPassword([]byte(pwd1), bcrypt.DefaultCost); err != nil {
		err = e
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
	Created          time.Time
	Modified         time.Time
	UID              string
	PWD              *password
	Country          string
	Postcode         string
	FirstName        string
	LastName         string
	ActivationCode   string
	PasswordCode     *passwordCode
	EmailAddressCode string
	NewUID           string
}

func (a *Account) FullName() (name string) {
	name = ""
	if len(a.FirstName) > 0 {
		name = a.FirstName
	}
	if len(a.LastName) > 0 {
		if len(name) > 0 {
			name = name + " "
		}
		name = name + a.LastName
	}
	if len(name) == 0 {
		name = a.UID
	}
	return
}

func (a *Account) Name() (name string) {
	if len(a.FirstName) > 0 {
		name = a.FirstName
	} else if len(a.LastName) > 0 {
		name = a.LastName
	} else {
		name = a.UID
	}
	return
}

func (a *Account) ValidateFields() (err error) {
	if false ||
		len(a.Country) == 0 ||
		len(a.Postcode) == 0 ||
		false {
		err = ErrValidationFailed
	}
	return
}

func (a *Account) Save() (err error) {
	a.Modified = time.Now()
	_, err = db.InsertUpdate(ACCOUNT_TABLE, a)
	return
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

func (a *Account) CreatePasswordCode(timeout time.Duration) error {
	a.PasswordCode = &passwordCode{
		Expires: time.Now().Add(timeout),
		Value:   code(),
	}
	return a.Save()
}

func ClearPasswordCode(uid, code string) {
	if acc, _, _ := get(uid); acc != nil {
		if acc.PasswordCode.Value == code {
			acc.PasswordCode = nil
			acc.Save()
		}
	}
}

func (a *Account) ChangePassword(code, pwd1, pwd2 string) (err error, conflict bool) {
	if code != a.PasswordCode.Value {
		err, conflict = ErrPasswordCodeIncorrect, true
	} else if pwd, e := newPassword(pwd1, pwd2); e != nil {
		err, conflict = e, true
	} else {
		a.PasswordCode = nil
		a.PWD = pwd
		err = a.Save()
	}
	return
}

func (a *Account) CreateEmailAddressCode(newUID string) error {
	a.NewUID = newUID
	a.EmailAddressCode = code()
	return a.Save()
}

func (a *Account) ClearEmailAddressCode(code string) (err error) {
	if a.EmailAddressCode == code {
		a.NewUID = ""
		a.EmailAddressCode = ""
		err = a.Save()
	}
	return
}

func (a *Account) ChangeEmailAddress(code string) (err error, conflict bool) {
	if code != a.EmailAddressCode {
		err, conflict = ErrEmailAddressCodeIncorrect, true
	} else {
		a.UID = a.NewUID
		err, conflict = a.ClearEmailAddressCode(code), false
	}
	return
}

func (a *Account) Refresh() (current bool) {
	if saved, e, _ := get(a.UID); e == nil {
		current = a.PWD.Created.Equal(saved.PWD.Created) && (a.ValidateFields() == nil)
		*a = *saved
	}
	return
}

func code() string {
	return string(util.URLEncode(util.Random()))
}

func New(uid, pwd1, pwd2 string) (account *Account, err error, conflict bool) {
	uid = strings.ToLower(uid)
	if e, found := db.Get(ACCOUNT_TABLE, uid, new(Account)); e != nil {
		err = e
	} else if found {
		err, conflict = ErrEmailTaken, true
	} else if pwd, e := newPassword(pwd1, pwd2); e != nil {
		err, conflict = e, true
	} else {
		acc := &Account{
			Created:        time.Now(),
			UID:            uid,
			PWD:            pwd,
			ActivationCode: code(),
		}
		if err = acc.Save(); err == nil {
			account = acc
		}
	}
	return
}

func get(uid string) (account *Account, err error, conflict bool) {
	acc := new(Account)
	if e, found := db.Get(ACCOUNT_TABLE, strings.ToLower(uid), acc); e != nil {
		err = e
	} else if !found {
		err, conflict = ErrInvalidCredentials, true
	} else {
		account = acc
	}
	return
}

func Activate(uid string, code string) (account *Account, err error, conflict bool) {
	if acc, e, c := get(uid); e != nil {
		err, conflict = e, c
	} else if e := acc.activate(code); e != nil {
		err, conflict = e, true
	} else if e := acc.Save(); e != nil {
		err = e
	} else {
		account = acc
	}
	return
}

func Get(uid, pwd string) (account *Account, err error, conflict bool) {
	if acc, e, c := get(uid); e != nil {
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

func GetInsecure(uid string) (account *Account, err error, conflict bool) {
	return get(uid)
}
