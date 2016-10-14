package routes

import (
	"github.com/wscherphof/essix/router"
	"github.com/wscherphof/essix/secure"
	"github.com/wscherphof/essix/template"
)

func init() {
	router.GET("/", secure.IfHandle(
		template.Handle("essix", "Home", "Home-LoggedIn"),
		template.Handle("essix", "Home", "Home-LoggedOut")))
}
