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

func Domain() string {
	// Die without a domain
	return env.Get("DOMAIN")
}
