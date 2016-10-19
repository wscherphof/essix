package routes

import (
	"github.com/wscherphof/essix/router"
	"github.com/wscherphof/essix/secure"
	"<routes/example>"
)

func init() {
	router.GET("/profile", secure.Handle(example.ProfileForm))
	router.PUT("/profile", secure.Handle(example.Profile))
}
