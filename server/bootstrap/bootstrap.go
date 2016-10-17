/*
Package bootstrap bootstraps Essix.

It reads the DOMAIN environment variable, and on initialisation connects to the
entity database, using the DB_NAME and DB_ADDRESS environment variables, or
their default values.
*/
package bootstrap

import (
	"github.com/wscherphof/entity"
	"github.com/wscherphof/env"
	"log"
)

func init() {
	var db = env.Get("DB_NAME", "essix")
	var address = env.Get("DB_ADDRESS", "db1")
	if err := entity.Connect(db, address); err != nil {
		log.Fatalf("ERROR: connecting to DB %s@%s failed. %T %s", db, address, err, err)
	} else {
		log.Printf("INFO: connected to DB %s@%s", db, address)
	}
}

/*
Domain returns the DOMAIN environment variable. A fatal error is triggered if
it's unset.
*/
func Domain() string {
	// Die without a domain
	return env.Get("DOMAIN")
}
