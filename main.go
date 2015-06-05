package main

import (
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/wscherphof/expeertise/captcha"
	"github.com/wscherphof/expeertise/env"
	"github.com/wscherphof/expeertise/router"
	"github.com/wscherphof/expeertise/secure"
	"log"
	"net/http"
	"os"
)

var (
	applicationAddress = env.Get("HTTP_HOST") + env.Get("HTTPS_PORT")
	redirectAddress    = env.Get("HTTP_HOST") + env.Get("HTTP_PORT")
)

func main() {
	// Redirect server serves http & redirects everything to the application server's https address
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+applicationAddress+r.URL.Path, http.StatusMovedPermanently)
		})
		log.Println("INFO: Redirect server    @", redirectAddress)
		log.Fatal(http.ListenAndServe(redirectAddress, nil))
	}()

	router.GET("/", secure.IfSecureHandle(
		router.Template("home", "home_loggedin", nil),
		router.Template("home", "home_loggedout", nil)))

	router.Router.Handler("GET", "/captcha/*filepath", captcha.Server)
	router.Router.ServeFiles("/static/*filepath", http.Dir("./static"))

	log.Println("INFO: Application server @", applicationAddress)
	log.Fatal(http.ListenAndServeTLS(applicationAddress, "cert.pem", "key.pem",
		context.ClearHandler(
			handlers.HTTPMethodOverrideHandler(
				handlers.CompressHandler(
					handlers.CombinedLoggingHandler(os.Stdout,
						router.Router))))))
}
