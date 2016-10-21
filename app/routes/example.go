package routes

import (
	"github.com/wscherphof/secure"
	"<routes/example>"
)

func init() {
	router.GET("/profile", secure.Handle(example.ProfileForm))
	router.PUT("/profile", secure.Handle(example.Profile))
}
