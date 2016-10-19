package secure

import (
	"encoding/gob"
	"github.com/gorilla/securecookie"
	"net/http"
	"strings"
	"time"
)

type RequestToken string

const requestTokenName = "4f3a0292-59c5-488a-b3fb-6e503c929331"

func NewRequestToken(r *http.Request, opt_path ...string) (token string, err error) {
	path := r.URL.Path
	if len(opt_path) == 1 {
		path = opt_path[0]
	}
	data := &RequestTokenData{
		IP:        IP(r),
		Timestamp: time.Now(),
		Path:      path,
	}
	return securecookie.EncodeMulti(requestTokenName, data, requestTokenCodecs...)
}

func (r RequestToken) Read(dst *RequestTokenData) (err error) {
	return securecookie.DecodeMulti(requestTokenName, string(r), dst, requestTokenCodecs...)
}

type RequestTokenData struct {
	IP        string
	Path      string
	Timestamp time.Time
}

func IP(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func init() {
	gob.Register(RequestTokenData{})
}
