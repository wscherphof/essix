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

func main() {
	router.GET("/", secure.IfSecureHandle(
		router.Template("home", "home_loggedin", nil),
		router.Template("home", "home_loggedout", nil)))

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+env.Get("HTTP_HOST")+env.Get("HTTPS_PORT")+r.URL.Path, http.StatusMovedPermanently)
		})
		address := env.Get("HTTP_HOST") + env.Get("HTTP_PORT")
		log.Println("INFO: Redirect server @", address)
		log.Fatal(http.ListenAndServe(address, nil))
	}()

	go func() {
		http.Handle("/static/", http.FileServer(http.Dir("")))
		http.Handle("/captcha/", captcha.Server)
		address := env.Get("HTTP_HOST") + env.Get("EXPRESS_PORT")
		log.Println("INFO: Express server  @", address)
		log.Fatal(http.ListenAndServeTLS(address, "cert.pem", "key.pem",
			handlers.CompressHandler(
				handlers.CombinedLoggingHandler(os.Stdout,
					http.DefaultServeMux))))
	}()

	address := env.Get("HTTP_HOST") + env.Get("HTTPS_PORT")
	log.Println("INFO: HTTPS server    @", address)
	log.Fatal(http.ListenAndServeTLS(address, "cert.pem", "key.pem",
		context.ClearHandler(
			secure.AuthenticationHandler(
				handlers.HTTPMethodOverrideHandler(
					handlers.CompressHandler(
						handlers.CombinedLoggingHandler(os.Stdout,
							router.Router)))))))
}
