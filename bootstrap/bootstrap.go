/*
Package bootstrap bootstraps Essix.

It reads the DOMAIN environment variable, and on initialisation connects to the
entity database, and initialises the messages data.
*/
package bootstrap

import (
	"github.com/wscherphof/entity"
	"github.com/wscherphof/env"
	"github.com/wscherphof/essix/messages"
)

func init() {
	entity.Connect(env.Get("DB_NAME", "essix"), env.Get("DB_ADDRESS", "db1"))
	messages.Init()
}

/*
Domain returns the DOMAIN environment variable. A fatal error is triggered if
it's unset.
*/
func Domain() string {
	// Die without a domain
	return env.Get("DOMAIN")
}
