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

	router.Router.Handler("GET", "/captcha/*filepath", captcha.Server)
	router.Router.ServeFiles("/static/*filepath", http.Dir("./static"))

	go func() {
		address := env.Get("HTTP_HOST") + env.Get("HTTP_PORT")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+address+r.URL.Path, http.StatusMovedPermanently)
		})
		log.Println("INFO: HTTP  server @", address)
		log.Fatal(http.ListenAndServe(address,
			handlers.CombinedLoggingHandler(os.Stdout,
				http.DefaultServeMux)))
	}()

	address := env.Get("HTTP_HOST") + env.Get("HTTPS_PORT")
	log.Println("INFO: HTTPS server @", address)
	log.Fatal(http.ListenAndServeTLS(address, "cert.pem", "key.pem",
		context.ClearHandler(
			secure.AuthenticationHandler(
				handlers.HTTPMethodOverrideHandler(
					handlers.CompressHandler(
						handlers.CombinedLoggingHandler(os.Stdout,
							router.Router)))))))
}
