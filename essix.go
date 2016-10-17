/*
Package essix runs an essential simple secure stable scalable stateless server.

See github.com/wscherphof/s6app for the canonocal usage example.

Essix relies on RethinkDB as its distributed database.
See https://wscherphof.wordpress.com/2016/09/13/rethink-swarm-mode/ for an easy
way to manage RethinkDB on Docker Swarm Mode, be it in your local development
environment, or in a multi-continent cloud setup.

Essix is designed to run on Swarm Mode in the same way, but isn't specifically
bound to that setup.

The `essix` script assumes the Swarm Mode setup though, and makes it a oneliner
to compile your app's sources, and build a Docker image, and/or run the image as
a service on the swarm.

Examples:
	$ essix -e DOMAIN=dev.wscherphof.nl build ../s6app wscherphof/s6app:0.2 s6app dev
Builds a 0.2 version of the wscherphof/s6app image on swarm dev's
nodes, and runs the s6app service on dev, with the DOMAIN variable
set to the given value.
	$ essix build ../s6app wscherphof/s6app:0.2 s6app
Locally builds a 0.2 version of the wscherphof/s6app image,
and pushes it to the repository.
	$ essix run wscherphof/s6app:0.2 s6app prd
Runs the s6app service on swarm prd, using the specified image,
which is downloaded from the repository, if not found locally.
*/
package essix

import (
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/wscherphof/essix/bootstrap"
	"github.com/wscherphof/essix/messages"
	"github.com/wscherphof/essix/router"
	"github.com/wscherphof/essix/routes"
	"log"
	"net/http"
	"os"
)

var (
	domain = bootstrap.Domain()
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
/resources/certificates.

See <future blogpost> for an easy way to obtain a certificate for your server.

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
