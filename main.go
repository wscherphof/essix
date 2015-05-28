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

const (
  HTTP_HOST  string = "localhost"
  HTTP_PORT  string = ":9090"
  HTTPS_PORT string = ":10443"
  DB_HOST string = "localhost"
  DB_PORT string = ":28015"
  DB_NAME string = "expeertise"
)

func main () {
  db.Init(DB_HOST + DB_PORT, DB_NAME)
  config.Init()
  secure.Init()
  model.Init()
  captcha.Init()
  DefineMessages()
  router := httprouter.New()

  router.GET    ("/", util.Template("home", "", nil))
  
  // TODO: sign up w/ just email & pwd; then on first login, ask further details
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
      http.Redirect(w, r, "https://" + HTTP_HOST + HTTPS_PORT + r.URL.Path, http.StatusMovedPermanently)
    })
    log.Fatal(http.ListenAndServe(HTTP_HOST + HTTP_PORT,
      handlers.CombinedLoggingHandler(os.Stdout,
    http.DefaultServeMux)))
  }()

  log.Fatal(http.ListenAndServeTLS(HTTP_HOST + HTTPS_PORT, "cert.pem", "key.pem",
    handlers.CombinedLoggingHandler(os.Stdout, 
    handlers.HTTPMethodOverrideHandler(
    handlers.CompressHandler(
    context.ClearHandler(
  router))))))
}
