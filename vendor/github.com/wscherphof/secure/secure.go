/*
Package secure manages client side session tokens for stateless web
applications.

Tokens are stored as an http cookie. An encrypted connection (https) is
required.

Call 'Configure()' once to provide the information for the package to operate,
including the type of the authentication data that will be used. The actual
configuration parameters are stored in a 'Config' type struct, which can be
synced with an external database, through the 'DB' interface.

Once configured, call 'Authentication()' to retrieve the data from the token.
It will redirect to a login page if no valid token is present (unless the
'optional' argument was 'true'). 'LogIn()' creates a new token, stores the
provided data in it, and redirects back to the page that required the
authentication.
'Update()' updates the authentication data in the current token. 'LogOut()'
deletes the token cookie.

You'll probably want to wrap 'Authentication()' in a function that converts the
'interface{}' result to the type that you use for the token data.
*/
package secure

import (
	"encoding/gob"
	"errors"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"time"
)

var (

	// ErrTokenNotSaved is returned by LogIn() if it couldn't set the cookie.
	ErrTokenNotSaved = errors.New("secure: failed to save the session token")

	// ErrNoTLS is returned by LogIn() if the connection isn't encrypted
	// (https)
	ErrNoTLS = errors.New("secure: logging in requires an encrypted conection")
)

const (
	authKeyLen = 32
	encrKeyLen = 32
)

const (
	tokenName      = "authtoken"
	recordField    = "ddf77ee1-6a23-4980-8edc-ff4139e98f22"
	createdField   = "45595a0b-7756-428e-bae0-5f7ded324e92"
	validatedField = "fe6f1315-9aa1-4083-89a0-dcb6c198654b"
	returnField    = "eb8cacdd-d65f-441e-a63d-e4da69c2badc"
)

// Config holds the package's configuration parameters.
// Can be synced with an external database, through the DB interface.
type Config struct {

	// LogInPath is the URL where Authentication() redirects to; a log in form
	// should be served here.
	// Default value is "/session".
	LogInPath string

	// LogOutPath is the URL where LogOut() redirects to.
	// Default value is "/".
	LogOutPath string

	// TimeOut is when a token expires (time after LogIn())
	// Default value is 6 * 30 days.
	TimeOut time.Duration

	// SyncInterval is how often the configuration is synced with an external
	// database. SyncInterval also determines whether it's time to have the
	// token data checked by the Validate function.
	// Default value is 5 minutes.
	SyncInterval time.Duration

	// KeyPairs are 4 32-long byte arrays (two pairs of an authentication key
	// and an encryption key); the 2nd pair is used for key rotation.
	// Default value is newly generated keys.
	// Keys get rotated on the first sync cycle after a TimeOut interval -
	// new tokens use the new keys; existing tokens continue to use the old
	// keys.
	KeyPairs [][]byte

	// TimeStamp is when the latest key pair was generated.
	TimeStamp time.Time

	RequestTokenKeyPairs [][]byte

	RequestTokenTimeStamp time.Time
}

var config = &Config{
	LogInPath:    "/session",
	LogOutPath:   "/",
	TimeOut:      6 * 30 * 24 * time.Hour,
	SyncInterval: 5 * time.Minute,
	KeyPairs: [][]byte{
		securecookie.GenerateRandomKey(authKeyLen),
		securecookie.GenerateRandomKey(encrKeyLen),
		securecookie.GenerateRandomKey(authKeyLen),
		securecookie.GenerateRandomKey(encrKeyLen),
	},
	TimeStamp: time.Now(),
	RequestTokenKeyPairs: [][]byte{
		securecookie.GenerateRandomKey(authKeyLen),
		securecookie.GenerateRandomKey(encrKeyLen),
		securecookie.GenerateRandomKey(authKeyLen),
		securecookie.GenerateRandomKey(encrKeyLen),
	},
	RequestTokenTimeStamp: time.Now(),
}

var (
	store              *sessions.CookieStore
	requestTokenCodecs []securecookie.Codec
)

func setKeys() {
	store = sessions.NewCookieStore(config.KeyPairs...)
	store.Options = &sessions.Options{
		MaxAge: int(config.TimeOut / time.Second),
		Secure: true,
		Path:   "/",
	}
	requestTokenCodecs = securecookie.CodecsFromPairs(config.RequestTokenKeyPairs...)
}

// DB is the interface to implement for syncing the configuration parameters.
//
// Syncing is executed every config.SyncInterval. If parameter values are
// changed in the database, the new values get synced to all servers that run
// the application.
type DB interface {

	// Fetch fetches a Config instance from the database.
	Fetch(dst *Config) error

	// Upsert inserts a Config instance into the database if none is present
	// on Configure(). Upsert updates the KeyPairs and TimeStamp values on key
	// rotation time.
	Upsert(src *Config) error
}

// Validate is the type of the function passed to Configure(), that gets called
// to have the application test whether the token data is still valid (e.g. to
// prevent continued access with a token that was created with an old password)
//
// 'src' is the authentication data from the token.
//
// 'dst' is the fresh authentication data to replace the token data with.
//
// 'valid' is whether the old data was good enough to keep the current token.
//
// Default implementation always returns the token data as is, and true, which
// is significantly insecure.
//
// Each successful validation stores a timestamp in
// the cookie. Validate is called on Authentication, if the time since the
// validation timestamp > config.SyncInterval
type Validate func(src interface{}) (dst interface{}, valid bool)

var validate = func(src interface{}) (dst interface{}, valid bool) {
	return src, true
}

// Configure configures the package and must be called once before calling any
// other function in this package.
//
// 'record' is an arbitrary (can be empty) instance of the type of the
// authentication data that will be passed to Login() to store in the token.
// It's needed to get its type registered with the serialisation package used
// (encoding/gob).
//
// 'db' is the implementation of the DB interface to sync the configuration
// parameters, or nil, in which case keys will not be rotated.
//
// 'validate' is the function that regularly verifies the token data, or nil,
// which would pose a significant security risk.
//
// 'optionalConfig' is the Config instance to start with. If omitted, the config
// from the db or the default config is used.
//
// Early experiments can skip the call to Configure(), and use a string or an
// int for the authentication data.
func Configure(record interface{}, db DB, validateFunc Validate, optionalConfig ...*Config) {
	gob.Register(record)
	gob.Register(time.Now())
	if len(optionalConfig) > 0 {
		opt := optionalConfig[0]
		if len(opt.LogInPath) > 0 {
			config.LogInPath = opt.LogInPath
		}
		if len(opt.LogOutPath) > 0 {
			config.LogOutPath = opt.LogOutPath
		}
		if opt.TimeOut > 0 {
			config.TimeOut = opt.TimeOut
		}
		if opt.SyncInterval > 0 {
			config.SyncInterval = opt.SyncInterval
		}
		if len(opt.KeyPairs) == 4 {
			config.KeyPairs = opt.KeyPairs
		}
		if !opt.TimeStamp.IsZero() {
			config.TimeStamp = opt.TimeStamp
		}
		if len(opt.RequestTokenKeyPairs) == 4 {
			config.RequestTokenKeyPairs = opt.RequestTokenKeyPairs
		}
		if !opt.RequestTokenTimeStamp.IsZero() {
			config.RequestTokenTimeStamp = opt.RequestTokenTimeStamp
		}
	}
	setKeys()
	if db != nil {
		go func() {
			for {
				sync(db)
				time.Sleep(config.SyncInterval)
			}
		}()
	}
	if validateFunc != nil {
		validate = validateFunc
	}
}

func sync(db DB) {
	dbConfig := new(Config)
	if err := db.Fetch(dbConfig); err != nil {
		// Upload current (default) config to DB if there wasn't any
		db.Upsert(config)
	} else {
		// Replace current config with the one from DB
		config = dbConfig
		// Rotate keys if timed out
		if time.Now().Sub(config.TimeStamp) > config.TimeOut {
			rotateConfig := new(Config)
			*rotateConfig = *config
			rotateConfig.KeyPairs = [][]byte{
				securecookie.GenerateRandomKey(authKeyLen),
				securecookie.GenerateRandomKey(encrKeyLen),
				config.KeyPairs[0],
				config.KeyPairs[1],
			}
			rotateConfig.TimeStamp = time.Now()
			if err := db.Upsert(rotateConfig); err != nil {
				config = rotateConfig
				log.Println("INFO: Security keys rotated")
			}
		}
		// Rotate RequestToken keys if timed out
		if time.Now().Sub(config.RequestTokenTimeStamp) > config.SyncInterval {
			rotateConfig := new(Config)
			*rotateConfig = *config
			rotateConfig.RequestTokenKeyPairs = [][]byte{
				securecookie.GenerateRandomKey(authKeyLen),
				securecookie.GenerateRandomKey(encrKeyLen),
				config.RequestTokenKeyPairs[0],
				config.RequestTokenKeyPairs[1],
			}
			rotateConfig.RequestTokenTimeStamp = time.Now()
			if err := db.Upsert(rotateConfig); err != nil {
				config = rotateConfig
				log.Println("INFO: RequestToken keys rotated")
			}
		}
		setKeys()
	}
}
