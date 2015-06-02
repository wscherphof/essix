package main

import (
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/wscherphof/expeertise/captcha"
	"github.com/wscherphof/expeertise/router"
	"github.com/wscherphof/expeertise/secure"
	"log"
	"net/http"
	"os"
)

const (
	// TODO: flags/envvars
	HTTP_HOST  = "localhost"
	HTTP_PORT  = ":9090"
	HTTPS_PORT = ":10443"
)

func main() {
	router.GET("/", secure.IfSecureHandle(
		router.Template("home", "home_loggedin", nil),
		router.Template("home", "home_loggedout", nil)))

	router.Router.Handler("GET", "/captcha/*filepath", captcha.Server)
	router.Router.ServeFiles("/static/*filepath", http.Dir("./static"))

	go func() {
		address := HTTP_HOST + HTTP_PORT
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+address+r.URL.Path, http.StatusMovedPermanently)
		})
		log.Println("INFO: HTTP  server @", address)
		log.Fatal(http.ListenAndServe(address,
			handlers.CombinedLoggingHandler(os.Stdout,
				http.DefaultServeMux)))
	}()

	address := HTTP_HOST + HTTPS_PORT
	log.Println("INFO: HTTPS server @", address)
	log.Fatal(http.ListenAndServeTLS(address, "cert.pem", "key.pem",
		context.ClearHandler(
			secure.AuthenticationHandler(
				handlers.HTTPMethodOverrideHandler(
					handlers.CompressHandler(
						handlers.CombinedLoggingHandler(os.Stdout,
							router.Router)))))))
}
