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

type FormToken struct {
	IP        string
	Path      string
	Timestamp time.Time
}

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

func (f *FormToken) String() (s string) {
	var err error
	if s, err = securecookie.EncodeMulti(formTokenName, f, requestTokenCodecs...); err != nil {
		log.Panicln("ERROR: encoding form token failed", err)
	}
	return
}

func (f *FormToken) Parse(s string) error {
	return securecookie.DecodeMulti(formTokenName, s, f, requestTokenCodecs...)
}

func ValidateFormToken(r *http.Request, opt_name ...string) (err error) {
	name := "_formtoken"
	if len(opt_name) == 1 {
		name = opt_name[0]
	}
	this, that := NewFormToken(r), new(FormToken)
	if err = that.Parse(r.FormValue(name)); err == nil {
		referer, _ := url.Parse(r.Referer())
		if that.IP != this.IP || (that.Path != this.Path && that.Path != referer.Path) {
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
