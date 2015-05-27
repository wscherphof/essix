package main

import (
  "net/http"
  "log"
  "os"
  "github.com/julienschmidt/httprouter"
  "github.com/gorilla/handlers"
  "github.com/gorilla/context"
  "github.com/wscherphof/expeertise/db"
  "github.com/wscherphof/expeertise/config"
  "github.com/wscherphof/expeertise/secure"
  "github.com/wscherphof/expeertise/model"
  "github.com/wscherphof/expeertise/captcha"
  "github.com/wscherphof/expeertise/util"
)

const HOST string = "localhost"
const HTTP_PORT string = ":9090"
const HTTPS_PORT string = ":10443"

func main () {
  db.Init("localhost:28015", "expeertise")
  config.Init()
  secure.Init()
  model.Init()
  captcha.Init()
  DefineMessages()
  router := httprouter.New()

  router.GET    ("/", util.Template("home", "", nil))
  
  // TODO: https
  
  router.GET    ("/account", secure.SignUpForm)
  router.POST   ("/account", secure.SignUp)

  router.GET    ("/account/activation",      secure.ActivateForm)
  router.GET    ("/account/activation/",     secure.ActivateForm)
  router.GET    ("/account/activation/:uid", secure.ActivateForm)
  router.PUT    ("/account/activation",      secure.Activate)
  
  router.GET    ("/account/activationcode",      secure.ActivationCodeForm)
  router.GET    ("/account/activationcode/",     secure.ActivationCodeForm)
  router.GET    ("/account/activationcode/:uid", secure.ActivationCodeForm)
  router.POST   ("/account/activationcode",      secure.ActivationCode)
  
  router.GET    ("/account/passwordcode",      secure.PasswordCodeForm)
  router.GET    ("/account/passwordcode/",     secure.PasswordCodeForm)
  router.GET    ("/account/passwordcode/:uid", secure.PasswordCodeForm)
  router.POST   ("/account/passwordcode",      secure.PasswordCode)
  
  router.GET    ("/account/password/:uid", secure.PasswordForm)
  router.PUT    ("/account/password",      secure.ChangePassword)
  
  router.GET    ("/session", secure.LogInForm)
  router.POST   ("/session", secure.LogIn)
  router.DELETE ("/session", secure.LogOut)
  
  // TODO:
  // router.GET    ("/account/instance/:uid", secure.Authenticate(secure.View))
  // router.PUT    ("/account/instance/:uid", secure.Authenticate(secure.Edit))
  // router.DELETE ("/account/instance/:uid", secure.Authenticate(secure.Terminate))

  router.GET    ("/protected", secure.Authenticate(Protected))
  
  router.Handler("GET", "/captcha/*filepath", captcha.Server)
  router.ServeFiles("/static/*filepath", http.Dir("./static"))

  go func(){
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      http.Redirect(w, r, "https://" + HOST + HTTPS_PORT + r.URL.Path, http.StatusMovedPermanently)
    })
    log.Fatal(http.ListenAndServe(HTTP_PORT,
      handlers.CombinedLoggingHandler(os.Stdout,
    http.DefaultServeMux)))
  }()

  log.Fatal(http.ListenAndServeTLS(HTTPS_PORT, "cert.pem", "key.pem",
    context.ClearHandler(
    handlers.HTTPMethodOverrideHandler(
    handlers.CompressHandler(
    handlers.CombinedLoggingHandler(os.Stdout, 
  router))))))

}
