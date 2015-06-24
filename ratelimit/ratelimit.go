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
	table = "ratelimit"
)

type path string

type client struct {
	IP       string
	Clear    time.Time
	Requests map[path]time.Time
}

func getClient(ip string) (c *client) {
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

func (c *client) save() (err error) {
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
	Path      path
	Timestamp time.Time
}

func NewToken(r *http.Request) (string, error) {
	return secure.NewRequestToken(&token{
		IP:        ip(r),
		Path:      path(r.URL.Path),
		Timestamp: time.Now(),
	})
}

func ip(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func Handle(seconds int, handle router.ErrorHandle) router.ErrorHandle {
	window := time.Duration(seconds) * time.Second
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
		t, ip, p := new(token), ip(r), path(r.URL.Path)
		if rate := r.FormValue("_rate"); rate == "" {
			err = router.NewError(ErrInvalidRequest)
			err.Conflict = true
			log.Printf("ATTACK: rate limit token missing %#v", r)
		} else if e := secure.RequestToken(rate).Read(t); e != nil {
			err = router.NewError(ErrInvalidRequest)
			err.Conflict = true
			log.Printf("ATTACK: rate limit token unreadable %#v", r)
		} else if t.IP != ip {
			err = router.NewError(ErrInvalidRequest)
			err.Conflict = true
			log.Printf("ATTACK: rate limit token invalid address: %v, not %v", t.IP, ip)
		} else if t.Path != p {
			err = router.NewError(ErrInvalidRequest)
			err.Conflict = true
			log.Printf("ATTACK: rate limit token invalid path: %v, not %v", t.Path, p)
		} else if c := getClient(ip); c.Requests[p].After(time.Now().Add(-window)) {
			err = router.NewError(ErrTooManyRequests)
			err.Conflict = true
		} else if c.Requests[p].After(t.Timestamp) {
			err = router.NewError(ErrTokenExpired)
			err.Conflict = true
		} else {
			c.Requests[p] = time.Now()
			c.Clear = time.Now().Add(window)
			if e := c.save(); e != nil {
				log.Printf("WARNING: error saving to table %v", table)
			}
			err = handle(w, r, ps)
		}
		return
	}
}

// TODO: clear job
