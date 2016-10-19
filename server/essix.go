/*
Package server runs the Essix server.
*/
package server

import (
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/wscherphof/essix/messages"
	"github.com/wscherphof/essix/router"
	"github.com/wscherphof/essix/routes"
	"github.com/wscherphof/env"
	"log"
	"net/http"
	"os"
)

var (
	// Die without a domain
	domain = env.Get("DOMAIN")
)

func init() {
	messages.Init()
	routes.Init()
}

/*
Run runs the application server. HTTP traffic on port 80 is redirected to HTTPS
on port 443.

Set the DOMAIN environment variable to the domain name to serve
HTTPS for; the certificates <DOMAIN>.crt & <DOMAIN>.key are expected in
/resources/certificates. `essix cert` generates certificates.

Set the DB_NAME & DB_ADDRESS environment variables for the RethinkDB connection.
Defaults: DB_NAME=essix, DB_ADDRESS=db1.
*/
func Run() {
	// Redirect http to https
	go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		url.Host = r.Host
		url.Scheme = "https"
		http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
	}))

	// Serve files in /static
	router.Router.ServeFiles("/static/*filepath", http.Dir("/resources/static"))

	log.Println("INFO: starting secure application server for " + domain)
	// Use the domain's proper certificates
	log.Fatal(http.ListenAndServeTLS(":443", "/resources/certificates/"+domain+".crt", "/resources/certificates/"+domain+".key",
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
}
