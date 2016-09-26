package main

import (
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/wscherphof/expeertise/env"
	"github.com/wscherphof/expeertise/router"
	"github.com/wscherphof/expeertise/secure"
	"log"
	"net/http"
	"os"
)

func init() {
	go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil || r.Host == "" {
			http.Error(w, "not found", http.StatusNotFound)
		}

		u := r.URL
		u.Host = r.Host
		u.Scheme = "https"
		http.Redirect(w, r, u.String(), http.StatusFound)
	}))
}

func main() {
	router.GET("/", secure.IfSecureHandle(
		router.Template(".", "home", "home_loggedin", nil),
		router.Template(".", "home", "home_loggedout", nil)))

	router.Router.ServeFiles("/static/*filepath", http.Dir("./static"))

	domain := env.Get("DOMAIN")
	log.Println("INFO: starting application server for " + domain)
	log.Fatal(http.ListenAndServeTLS(":443", "/certificates/"+domain+".crt", "/certificates/"+domain+".key",
		context.ClearHandler(
			handlers.HTTPMethodOverrideHandler(
				handlers.CompressHandler(
					handlers.CombinedLoggingHandler(os.Stdout,
						router.Router))))))
}
