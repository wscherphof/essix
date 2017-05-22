/*
Package secure manages authentication cookies for stateless web applications,
and form tokens for CSRF protection.

An encrypted connection (https) is required.

Call 'Configure()' once to provide the information for the package to operate,
including the type of the authentication data that will be used. The actual
configuration parameters are stored in a 'Config' type struct. The 'DB'
interface syncs the Config to an external database, and automatically rotates
security keys.

Once configured, call 'Authentication()' to retrieve the data from the cookie.
It will redirect to a login page if no valid cookie is present (unless the
'optional' argument was 'true'). 'LogIn()' creates a new cookie, stores the
provided data in it, and redirects back to the page that required the
authentication.
'Update()' updates the authentication data in the current cookie. 'LogOut()'
deletes the cookie.

You'll probably want to wrap 'Authentication()' in a function that converts the
'interface{}' result to the type that you use for the cookie data.
*/
package secure

import (
	"encoding/gob"
	"errors"
	"github.com/gorilla/securecookie"
	"log"
	"time"
)

var (

	// ErrTokenNotSaved is returned by LogIn() if it couldn't set the cookie.
	ErrTokenNotSaved = errors.New("secure: failed to save the session cookie")

	// ErrNoTLS is returned by LogIn() if the connection isn't encrypted
	// (https)
	ErrNoTLS = errors.New("secure: logging in requires an encrypted conection")
)

const (
	authKeyLen = 32
	encrKeyLen = 32
)

// Keys manages rotating cryptographic keys.
type Keys struct {

	// KeyPairs holds 3 key pairs: current, previous, and next.
	KeyPairs [][]byte

	// Start is when the current key pair became current.
	Start time.Time

	// TimeOut is how much time Start the key pairs should be rotated.
	TimeOut time.Duration
}

func (k *Keys) stale() bool {
	return time.Since(k.Start) >= k.TimeOut
}

func (k *Keys) rotate() (ret *Keys) {
	ret = &Keys{
		KeyPairs: [][]byte{
			k.KeyPairs[4],
			k.KeyPairs[5],
			k.KeyPairs[0],
			k.KeyPairs[1],
			securecookie.GenerateRandomKey(authKeyLen),
			securecookie.GenerateRandomKey(encrKeyLen),
		},
		TimeOut: k.TimeOut,
		Start:   time.Now(),
	}
	return ret
}

func (k *Keys) freshen() {
	if k.stale() {
		// Immediately stop encoding with the stale key, and start using the
		// commonly known next key as the current one for this time window.
		*k = *k.rotate()
		// Parallelly update the DB to establish a new commonly known next key
		// before the next time window starts.
		go sync()
	}
}

// Config holds the package's configuration parameters.
// Can be synced with an external database, through the DB interface.
type Config struct {

	// Session manages session cryptography.
	Session *Session

	// Token manages token cryptography.
	Token *Token

	// Locked prevents interfering updates when syncing.
	Locked bool
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
	// on Configure(). Upsert updates the CookieKeyPairs and CookieTimeStamp values on key
	// rotation time.
	Upsert(src *Config) error
}

func (c *Config) wait() {
	for c.Locked {
		time.Sleep(50 * time.Millisecond)
		if err := db.Fetch(c); err != nil {
			log.Println("WARNING: secure config wait: fetch failed:", err)
		}
	}
}

func (c *Config) lock() {
	if !c.Locked {
		c.Locked = true
		if err := db.Upsert(c); err != nil {
			c.Locked = false
			log.Panicln("ERROR: secure config lock failed:", err)
		}
	}
}

func (c *Config) unlock() {
	if c.Locked {
		c.Locked = false
		if err := db.Upsert(c); err != nil {
			c.Locked = true
			log.Panicln("ERROR: secure config unlock failed:", err)
		}
	}
}

func sync() {
	dbConfig := new(Config)
	if err := db.Fetch(dbConfig); err != nil {
		// Upload current (default) config to DB if there wasn't any
		if err := db.Upsert(config); err != nil {
			log.Panicln("ERROR: secure DB: saving default config failed:", err)
		}
	} else {
		// Replace current config with the one from DB
		config = dbConfig
		// Rotate keys if timed out
		config.wait()
		if config.Session.stale() {
			log.Println("INFO: secure DB: rotating session keys...")
			config.lock()
			config.Session.Keys = config.Session.rotate()
		}
		if config.Token.stale() {
			log.Println("INFO: secure DB: rotating token keys...")
			config.lock()
			config.Token.Keys = config.Token.rotate()
		}
		if config.Locked {
			if err := db.Upsert(config); err != nil {
				log.Panicln("ERROR: secure DB: key rotatation failed:", err)
			} else {
				log.Println("INFO: secure DB: keys rotated")
			}
			config.unlock()
		}
	}
}

// ValidateCookie is the type of the function passed to Configure(), that gets called
// to have the application test whether the cookie data is still valid (e.g. to
// prevent continued access with a cookie that was created with an old password)
//
// 'src' is the authentication data from the cookie.
//
// 'dst' is the fresh authentication data to replace the cookie data with.
//
// 'valid' is whether the old data was good enough to keep the current cookie.
//
// Default implementation always returns the cookie data as is, and true, which
// is significantly insecure.
//
// Each successful validation stores a timestamp in
// the cookie. ValidateCookie is called on Authentication, if the time since the
// validation timestamp > config.SyncInterval
type ValidateCookie func(src interface{}) (dst interface{}, valid bool)

var (
	db       DB
	validate ValidateCookie
	config   = &Config{
		Session: &Session{
			Keys: &Keys{
				KeyPairs: [][]byte{
					securecookie.GenerateRandomKey(authKeyLen),
					securecookie.GenerateRandomKey(encrKeyLen),
					securecookie.GenerateRandomKey(authKeyLen),
					securecookie.GenerateRandomKey(encrKeyLen),
					securecookie.GenerateRandomKey(authKeyLen),
					securecookie.GenerateRandomKey(encrKeyLen),
				},
				Start:   time.Now(),
				TimeOut: 6 * 30 * 24 * time.Hour,
			},
			LogInPath:       "/session",
			LogOutPath:      "/",
			ValidateTimeOut: 5 * time.Minute,
		},
		Token: &Token{
			Keys: &Keys{
				KeyPairs: [][]byte{
					securecookie.GenerateRandomKey(authKeyLen),
					securecookie.GenerateRandomKey(encrKeyLen),
					securecookie.GenerateRandomKey(authKeyLen),
					securecookie.GenerateRandomKey(encrKeyLen),
					securecookie.GenerateRandomKey(authKeyLen),
					securecookie.GenerateRandomKey(encrKeyLen),
				},
				Start:   time.Now(),
				TimeOut: 5 * time.Minute,
			},
		},
	}
)

// Configure configures the package and must be called once before calling any
// other function in this package.
//
// 'record' is an arbitrary (can be empty) instance of the type of the
// authentication data that will be passed to Login() to store in the cookie.
// It's needed to get its type registered with the serialisation package used
// (encoding/gob).
//
// 'dbImpl' is the implementation of the DB interface to sync the configuration
// and rotate the keys.
//
// 'validate' is the function that regularly verifies the cookie data.
//
// 'opt_config' is the Config instance to start with. If omitted, the config
// from the db is used, or else the default config.
func Configure(record interface{}, dbImpl DB, validateFunc ValidateCookie, opt_config ...*Config) {
	gob.Register(record)
	gob.Register(time.Now())
	db = dbImpl
	validate = validateFunc
	if len(opt_config) == 1 {
		config = opt_config[0]
	}
	sync()
	go func() {
		time.Sleep(config.Session.TimeOut / 2)
		for {
			sync()
			// will the timeout get reread each iteration, to cater for updates?
			time.Sleep(config.Session.TimeOut)
		}
	}()
	go func() {
		time.Sleep(config.Token.TimeOut / 2)
		for {
			sync()
			// will the timeout get reread each iteration, to cater for updates?
			time.Sleep(config.Token.TimeOut)
		}
	}()
}
