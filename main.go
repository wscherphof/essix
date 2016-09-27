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

func main() {
	// Serve files in /static
	router.Router.ServeFiles("/static/*filepath", http.Dir("./rrresources/static"))

	// Template for home page, depending on login status
	router.GET("/", secure.IfSecureHandle(
		router.Template(".", "home", "home_loggedin", nil),
		router.Template(".", "home", "home_loggedout", nil)))

	domain := env.Default("DOMAIN", "dev.wscherphof.nl")
	log.Println("INFO: starting secure application server for " + domain)
	// Use the domain's proper certificates
	log.Fatal(http.ListenAndServeTLS(":443", "/rrresources/certificates/"+domain+".crt", "/rrresources/certificates/"+domain+".key",
		// Clear the context data created for the request, as per the "Important note" in https://godoc.org/github.com/gorilla/sessions
		context.ClearHandler(
			// Support PUT & DELTE through POST forms
			handlers.HTTPMethodOverrideHandler(
				// Zip responses
				handlers.CompressHandler(
					// Log request info in Apache Combined Log Format
					handlers.CombinedLoggingHandler(os.Stdout,
						// Use our routes
						router.Router))))))

	// Redirect http to https
	go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL
		u.Host = r.Host
		u.Scheme = "https"
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	}))
}
