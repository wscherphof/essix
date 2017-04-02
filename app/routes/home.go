package routes

import (
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/secure"
)

func init() {
	router.GET("/", secure.IfHandle(
		template.Handle("essix", "Home", "Home-LoggedIn"),
		template.Handle("essix", "Home", "Home-LoggedOut")))
}
