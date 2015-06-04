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
	"strings"
)

var (
	applicationAddress = env.Get("HTTP_HOST") + env.Get("HTTPS_PORT")
	expressAddress     = env.Get("HTTP_HOST") + env.Get("EXPRESS_PORT")
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

	// Express server serves static resources through https, without application handler overhead
	go func() {
		fileServer := http.FileServer(http.Dir(""))
		http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r)
			} else {
				fileServer.ServeHTTP(w, r)
			}
		})
		http.Handle("/captcha/", captcha.Server)
		log.Println("INFO: Express server     @", expressAddress)
		log.Fatal(http.ListenAndServeTLS(expressAddress, "cert.pem", "key.pem",
			handlers.CompressHandler(
				handlers.CombinedLoggingHandler(os.Stdout,
					http.DefaultServeMux))))
	}()

	// Home is the application's start page
	router.GET("/", secure.IfSecureHandle(
		router.Template("home", "home_loggedin", nil),
		router.Template("home", "home_loggedout", nil)))

	// Application server wraps and handles application routes
	log.Println("INFO: Application server @", applicationAddress)
	log.Fatal(http.ListenAndServeTLS(applicationAddress, "cert.pem", "key.pem",
		context.ClearHandler(
			handlers.HTTPMethodOverrideHandler(
				handlers.CompressHandler(
					handlers.CombinedLoggingHandler(os.Stdout,
						router.Router))))))
}
