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
func (f *FormToken) String() (s string) {
	var err error
	if s, err = securecookie.EncodeMulti(formTokenName, f, formTokenCodecs...); err != nil {
		log.Panicln("ERROR: encoding form token failed", err)
	}
	return
}

/*
Parse populates the data fields from an encrypted token string.
*/
func (f *FormToken) Parse(s string) error {
	return securecookie.DecodeMulti(formTokenName, s, f, formTokenCodecs...)
}

func ip(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func init() {
	gob.Register(FormToken{})
}
