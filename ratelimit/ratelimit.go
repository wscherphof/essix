package ratelimit

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/expeertise/db"
	"github.com/wscherphof/expeertise/router"
	"github.com/wscherphof/secure"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	ErrTooManyRequests = errors.New("429 Too Many Requests")
	ErrInvalidRequest  = errors.New("Invalid Request")
	ErrTokenExpired    = errors.New("Token Expired")
)

const (
	table        = "ratelimit"
	tokenTimeOut = time.Minute
)

type path string

type client struct {
	IP       string
	Clear    time.Time
	Requests map[path]time.Time
}

func GetClient(ip string) (c *client) {
	c = new(client)
	err, found := db.Get(table, ip, c)
	if err != nil {
		log.Printf("WARNING: error getting from table %v", table)
	}
	if !found {
		c.IP = ip
		c.Requests = make(map[path]time.Time)
	}
	return
}

func (c *client) Save() (err error) {
	_, err = db.InsertUpdate(table, c)
	return
}

func init() {
	if _, err := db.TableCreatePK(table, "IP"); err == nil {
		log.Println("INFO: table created:", table)
	}
	secure.RegisterRequestTokenData(token{})
}

type token struct {
	IP        string
	Timestamp time.Time
	// TODO: Path path
}

func NewToken(r *http.Request) (string, error) {
	return secure.NewRequestToken(&token{
		IP:        ip(r),
		Timestamp: time.Now(),
	})
}

func ip(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func Handle(seconds int, handle router.ErrorHandle) router.ErrorHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
		t := new(token)
		if rate := r.FormValue("_rate"); rate == "" {
			err = router.NewError(ErrInvalidRequest)
			err.Conflict = true
			log.Printf("ATTACK: rate limit token missing %#v", r)
		} else if e := secure.RequestToken(rate).Read(t); e != nil {
			err = router.NewError(ErrInvalidRequest)
			err.Conflict = true
			log.Printf("ATTACK: rate limit token unreadable %#v", r)
		} else if t.IP != ip(r) {
			err = router.NewError(ErrInvalidRequest)
			err.Conflict = true
			log.Printf("ATTACK: rate limit token invalid address: %v, not %v", t.IP, ip(r))
		} else if t.Timestamp.After(time.Now().Add(tokenTimeOut)) {
			err = router.NewError(ErrTokenExpired)
			err.Conflict = true
		}
		if err != nil {
			return err
		}

		c := GetClient(ip(r))
		p := path(r.URL.Path)
		if !c.Requests[p].IsZero() && c.Requests[p].After(time.Now().Add(time.Duration(-seconds)*time.Second)) {
			err = router.NewError(ErrTooManyRequests)
			err.Conflict = true
		} else if c.Requests[p].After(t.Timestamp) {
			err = router.NewError(ErrTokenExpired)
			err.Conflict = true
		} else {
			c.Requests[p] = time.Now()
			c.Clear = time.Now().Add(time.Duration(seconds) * time.Second)
			if e := c.Save(); e != nil {
				log.Printf("WARNING: error saving to table %v", table)
			}
			err = handle(w, r, ps)
		}
		return
	}
}

// TODO: clear job
