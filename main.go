package main

import (
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/wscherphof/expeertise/router"
	"github.com/wscherphof/expeertise/secure"
	"github.com/wscherphof/letsencrypt"
	"io"
	"log"
	"net/http"
	"os"
)

func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

func main() {

	router.GET("/", secure.IfSecureHandle(
		router.Template(".", "home", "home_loggedin", nil),
		router.Template(".", "home", "home_loggedout", nil)))

	router.Router.ServeFiles("/static/*filepath", http.Dir("./static"))

	if _, err := os.Stat("/appdata/letsencrypt.cache"); os.IsNotExist(err) {
		if err := copy("/letsencrypt.cache", "/appdata/letsencrypt.cache"); err != nil {
			log.Fatal("ERROR: ", err)
		}
	}

	var m letsencrypt.Manager
	if err := m.CacheFile("/appdata/letsencrypt.cache"); err != nil {
		log.Fatal(err)
	}

	log.Println("INFO: starting application server")
	log.Fatal(m.Serve(
		context.ClearHandler(
			handlers.HTTPMethodOverrideHandler(
				handlers.CompressHandler(
					handlers.CombinedLoggingHandler(os.Stdout,
						router.Router))))))
}
