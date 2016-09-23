package main

import (
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/wscherphof/expeertise/router"
	"github.com/wscherphof/expeertise/secure"
	"github.com/wscherphof/letsencrypt"
	"log"
	"net/http"
	"os"
)

func main() {

	router.GET("/", secure.IfSecureHandle(
		router.Template(".", "home", "home_loggedin", nil),
		router.Template(".", "home", "home_loggedout", nil)))

	router.Router.ServeFiles("/static/*filepath", http.Dir("./static"))

	var m letsencrypt.Manager
	if err := m.CacheFile("/appdata/letsencrypt.cache"); err != nil {
		log.Fatal(err)
	}

	log.Println("INFO: starting application server")
	log.Fatal(m.Serve(
		context.ClearHandler(
			handlers.HTTPMethodOverrideHandler(
				handlers.CompressHandler(
					handlers.CombinedLoggingHandler(os.Stdout,
						router.Router))))))
}
