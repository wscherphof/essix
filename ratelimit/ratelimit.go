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
)

const (
	table = "ratelimit"
)

func init() {
	if _, err := db.TableCreatePK(table, "IP"); err == nil {
		log.Println("INFO: table created:", table)
	}
	secure.RegisterRequestTokenData(token{})
}

type path string

type token struct {
	IP        string
	Path      path
	Timestamp time.Time
}

func ip(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func NewToken(r *http.Request) (string, error) {
	return secure.NewRequestToken(&token{
		IP:        ip(r),
		Path:      path(r.URL.Path),
		Timestamp: time.Now(),
	})
}

type requests map[path]time.Time

type client struct {
	IP       string
	Clear    time.Time
	Requests requests
}

func getClient(ip string) (c *client) {
	c = new(client)
	err, found := db.Get(table, ip, c)
	if err != nil {
		log.Printf("WARNING: error getting from table %v: %v", table, err)
	}
	if !found {
		c.IP = ip
		c.Requests = make(requests)
	}
	return
}

func (c *client) save() (err error) {
	_, err = db.InsertUpdate(table, c)
	return
}

func Handle(seconds int, handle router.ErrorHandle) router.ErrorHandle {
	window := time.Duration(seconds) * time.Second
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
		t, ip, p := new(token), ip(r), path(r.URL.Path)
		if rate := r.FormValue("_rate"); rate == "" {
			err = router.NewError(ErrInvalidRequest)
			err.Conflict = true
			log.Printf("SUSPICIOUS: rate limit token missing %v %v", ip, p)
		} else if e := secure.RequestToken(rate).Read(t); e != nil {
			err = router.NewError(ErrInvalidRequest)
			err.Conflict = true
			log.Printf("SUSPICIOUS: rate limit token unreadable %v %v", ip, p)
		} else if t.IP != ip {
			err = router.NewError(ErrInvalidRequest)
			err.Conflict = true
			log.Printf("SUSPICIOUS: rate limit token invalid address: %v, expected %v %v", t.IP, ip, p)
		} else if t.Path != p {
			err = router.NewError(ErrInvalidRequest)
			err.Conflict = true
			log.Printf("SUSPICIOUS: rate limit token invalid path: %v, token path %v, expected %v", ip, t.Path, p)
		} else if c := getClient(ip); c.Requests[p].After(t.Timestamp) {
			err = router.NewError(ErrInvalidRequest)
			err.Conflict = true
			log.Printf("SUSPICIOUS: rate limit token reuse: %v %v, token %v, previous request %v", ip, p, t.Timestamp, c.Requests[p])
		} else if c.Requests[p].After(time.Now().Add(-window)) {
			err = router.NewError(ErrTooManyRequests, "ratelimit", "toomanyrequests")
			err.Conflict = true
			err.Data = map[string]interface{}{
				"Window": window,
			}
		} else {
			c.Requests[p] = time.Now()
			clear := time.Now().Add(window)
			if clear.After(c.Clear) {
				c.Clear = clear
			}
			if e := c.save(); e != nil {
				log.Printf("WARNING: error saving to table %v: %v", table, e)
			}
			err = handle(w, r, ps)
		}
		return
	}
}

// TODO: clear job
