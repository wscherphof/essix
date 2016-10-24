package routes

import (
	"app/routes/example"
	"github.com/wscherphof/secure"
)

func init() {
	router.GET("/profile", secure.Handle(example.ProfileForm))
	router.PUT("/profile", secure.Handle(example.Profile))
}
