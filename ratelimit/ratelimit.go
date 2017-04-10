/*
Package ratelimit manages rate limits for urls.

It can be used instead of captchas, since breaking captchas is a commercial
service, yielding a captcha to be nothing more than a rate limit.

The rate limited requests for each client (ip address) are stored in the
database, and subsequent requests from the same client are denied if they come
within the set time window.

Entries of non-returning clients are cleared from the database regularly.
*/
package ratelimit

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/entity"
	"github.com/wscherphof/env"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/secure"
	"log"
	"net/http"
	"time"
)

var (
	ErrTooManyRequests = errors.New("ErrTooManyRequests")
	ErrInvalidRequest  = errors.New("ErrInvalidRequest")
	defaultLimit       = env.GetInt("RATELIMIT", 60)
)

const (
	clearInterval = time.Hour
)

type requests map[string]time.Time

type client struct {
	*entity.Base
	Clear    int64
	Requests requests
}

func initClient(opt_id ...string) (c *client) {
	c = &client{Base: &entity.Base{}}
	if len(opt_id) == 1 {
		c.ID = opt_id[0]
	}
	return
}

func init() {
	c := initClient()
	entity.Register(c).Index("Clear")
	go func() {
		for {
			time.Sleep(clearInterval)
			limit := time.Now().Unix()
			index := c.Index(c, "Clear")
			selection := index.Between(nil, true, limit, true)
			if deleted, err := selection.Delete(); err != nil {
				log.Printf("WARNING: rate limit clearing failed: %v", err)
			} else {
				log.Printf("INFO: %v rate limit records cleared", deleted)
			}
		}
	}()
}

func getClient(ip string) (c *client) {
	c = initClient(ip)
	if err, empty := c.Read(c); err != nil {
		if empty {
			c.Requests = make(requests, 20)
		} else {
			log.Printf("WARNING: error reading from ratelimit table: %v", err)
		}
	}
	return
}

/*
ratelimit.Handle returns a httprouter.Handle that denies a request if it's
repeated from the same client within the given number of seconds, or handles it,
and resets the timw window.

Default limit is read from the RATELIMIT environment variable. The default value
for RATELIMIT is 60 seconds.
*/
func Handle(handle httprouter.Handle, opt_seconds ...int) httprouter.Handle {
	seconds := defaultLimit
	if len(opt_seconds) == 1 {
		seconds = opt_seconds[0]
	}
	window := time.Duration(seconds) * time.Second
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		this, that := secure.NewFormToken(r), new(secure.FormToken)
		if tokenString := r.FormValue("_formtoken"); tokenString == "" {
			template.Error(w, r, ErrInvalidRequest, true)
			log.Printf("SUSPICIOUS: rate limit token missing %v %v", this.IP, this.Path)
		} else if e := that.Parse(tokenString); e != nil {
			template.Error(w, r, ErrInvalidRequest, true)
			log.Printf("SUSPICIOUS: rate limit token unreadable %v %v", this.IP, this.Path)
		} else if that.IP != this.IP {
			template.Error(w, r, ErrInvalidRequest, true)
			log.Printf("SUSPICIOUS: rate limit token invalid address: %v, expected %v %v", that.IP, this.IP, this.Path)
		} else if that.Path != this.Path {
			template.Error(w, r, ErrInvalidRequest, true)
			log.Printf("SUSPICIOUS: rate limit token invalid path: %v, token path %v, expected %v", this.IP, that.Path, this.Path)
		} else if c := getClient(this.IP); c.Requests[this.Path].After(that.Timestamp) {
			template.Error(w, r, ErrInvalidRequest, true)
			log.Printf("SUSPICIOUS: rate limit token reuse: %v %v, token %v, previous request %v", this.IP, this.Path, that.Timestamp, c.Requests[this.Path])
		} else if c.Requests[this.Path].After(time.Now().Add(-window)) {
			template.ErrorTail(w, r, ErrTooManyRequests, true, "ratelimit", "TooManyRequests-error-tail", map[string]interface{}{
				"window": window,
			})
		} else {
			c.Requests[this.Path] = time.Now()
			if clear := time.Now().Add(window).Unix(); clear > c.Clear {
				c.Clear = clear
			}
			if e := c.Update(c); e != nil {
				log.Printf("WARNING: error updating ratelimit record: %v", e)
			}
			handle(w, r, ps)
		}
		return
	}
}
