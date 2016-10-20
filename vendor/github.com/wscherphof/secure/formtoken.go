package secure

import (
	"encoding/gob"
	"errors"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
	"net/url"
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

/*
ValidateFormToken reads an encrypted token string from the request, and returns
nil if the token is valid.
The token is expected in r.FormValue(opt_name); default name is "_formtoken"
*/
func ValidateFormToken(r *http.Request, opt_name ...string) (err error) {
	name := "_formtoken"
	if len(opt_name) == 1 {
		name = opt_name[0]
	}
	this, that := NewFormToken(r), new(FormToken)
	if err = that.Parse(r.FormValue(name)); err == nil {
		referer, _ := url.Parse(r.Referer())
		if that.IP != this.IP || (that.Path != this.Path && that.Path != referer.Path) {
			// Timestamp not considered, since key rotation will outdate old tokens automatically
			err = errors.New("Form token invalid")
		}
	}
	if err != nil {
		log.Printf("WARNING: %s %s %s", err, this.IP, this.Path)
	}
	return
}

func ip(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func init() {
	gob.Register(FormToken{})
}
