package essix

import (
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/wscherphof/essix/bootstrap"
	"github.com/wscherphof/essix/router"
	"github.com/wscherphof/essix/secure"
	"github.com/wscherphof/essix/template"
	"log"
	"net/http"
	"os"
)

var domain = bootstrap.Domain()

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

	// Template for home page, depending on login status
	router.GET("/", secure.IfSecureHandle(
		template.Handle("essix", "Home", "Home-LoggedIn", nil),
		template.Handle("essix", "Home", "Home-LoggedOut", nil)))

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
