package essix

import (
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/wscherphof/essix/env"
	"github.com/wscherphof/essix/router"
	"github.com/wscherphof/essix/messages"
	"github.com/wscherphof/essix/secure"
	"log"
	"net/http"
	"os"
)

var domain string

func init() {
	domain = env.Get("DOMAIN")
	if domain == "" {
		log.Fatal("DOMAIN environment variable not set")
	}
	messages.Init()
}

func Run() {
	// Serve files in /static
	router.Router.ServeFiles("/static/*filepath", http.Dir("/resources/static"))

	// Template for home page, depending on login status
	router.Router.GET("/", secure.IfSecureHandle(
		router.Template("essix", "home", "home_loggedin", nil),
		router.Template("essix", "home", "home_loggedout", nil)))

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

	// Redirect http to https
	go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL
		u.Host = r.Host
		u.Scheme = "https"
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	}))
}
