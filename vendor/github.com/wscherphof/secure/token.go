package secure

import (
	"encoding/gob"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
	"strings"
	"time"
)

/*
A Token value is stored in the Config to manage token cryptography.
*/
type Token struct {

	// Keys encapsulates the rotating key data & functionality.
	*Keys

	_codecs []securecookie.Codec
}

func (t *Token) codecs() []securecookie.Codec {
	if len(t._codecs) == 0 {
		t._codecs = securecookie.CodecsFromPairs(t.KeyPairs...)
	}
	return t._codecs
}

func (t *Token) encode(name string, value interface{}) (s string) {
	t.freshen()
	var err error
	if s, err = securecookie.EncodeMulti(name, value, t.codecs()...); err != nil {
		log.Panicln("ERROR: encoding form token failed", err)
	}
	return
}

func (t *Token) decode(name string, value string, dst interface{}) error {
	return securecookie.DecodeMulti(name, value, dst, t.codecs()...)
}

/*
A FormToken is a secured identification of a request, suitable for protection
against cross site request forgery.
*/
type FormToken struct {

	// IP is the client's IP address.
	IP string

	// Path is the path on this server the token is valid for.
	Path string

	// Timestamp is when the token was created.
	Timestamp time.Time
}

/*
NewFormToken initialises a new FormToken.
The opt_path argument overrides the default r.URL.Path.
*/
func NewFormToken(r *http.Request, opt_path ...string) *FormToken {
	path := r.URL.Path
	if len(opt_path) == 1 {
		path = opt_path[0]
	}
	return &FormToken{
		IP:        ip(r),
		Timestamp: time.Now(),
		Path:      path,
	}
}

const formTokenName = "4f3a0292-59c5-488a-b3fb-6e503c929331"

/*
String returns the encrypted token string.
*/
func (f *FormToken) String() string {
	return tokenKeys.encode(formTokenName, f)
}

/*
Parse populates the data fields from an encrypted token string.
*/
func (f *FormToken) Parse(s string) error {
	return tokenKeys.decode(formTokenName, s, f)
}

func ip(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func init() {
	gob.Register(FormToken{})
}
