package secure

import (
	"encoding/gob"
	"github.com/gorilla/securecookie"
	"net/http"
	"strings"
	"time"
	"log"
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

func ip(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func init() {
	gob.Register(FormToken{})
}
