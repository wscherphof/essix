package main

import (
  "net/http"
  "log"
  "os"
  "github.com/gorilla/handlers"
  "github.com/gorilla/context"
  "github.com/wscherphof/expeertise/router"
  "github.com/wscherphof/expeertise/db"
  "github.com/wscherphof/expeertise/config"
  "github.com/wscherphof/expeertise/secure"
  "github.com/wscherphof/expeertise/model"
  "github.com/wscherphof/expeertise/captcha"
  "github.com/wscherphof/expeertise/util2"
)

const (
  // TODO: flags/envvars
  HTTP_HOST  = "localhost"
  HTTP_PORT  = ":9090"
  HTTPS_PORT = ":10443"
  DB_HOST    = "localhost"
  DB_PORT    = ":28015"
  DB_NAME    = "expeertise"
)

func main () {
  // TODO: init()
  db.Init(DB_HOST + DB_PORT, DB_NAME)
  config.Init()
  secure.Init()
  model.Init()
  captcha.Init()
  DefineMessages()

  // TODO: differentiate whether logged in
  router.GET    ("/", util2.Template("home", "", nil))
  
  // TODO: sign up w/ just email & pwd; then on first login, ask further details
  // TODO: change email address (only when logged in, but still w/ an email to the new address)
  router.GET    ("/account", secure.IfSecureHandle(secure.UpdateAccountForm, secure.SignUpForm))
  router.POST   ("/account", secure.SignUp)
  router.PUT    ("/account", secure.SecureHandle(secure.UpdateAccount))
  // TODO: router.DELETE ("/account", secure.Authenticate(secure.TerminateAccount))

  // TODO: router
  router.Router.GET    ("/session", secure.LogInForm)
  router.Router.POST   ("/session", secure.LogIn)
  router.Router.DELETE ("/session", secure.LogOut)

  router.Router.GET    ("/account/activation",      secure.ActivateForm)
  router.Router.GET    ("/account/activation/",     secure.ActivateForm)
  router.Router.GET    ("/account/activation/:uid", secure.ActivateForm)
  router.Router.PUT    ("/account/activation",      secure.Activate)
  
  router.Router.GET    ("/account/activationcode",      secure.ActivationCodeForm)
  router.Router.GET    ("/account/activationcode/",     secure.ActivationCodeForm)
  router.Router.GET    ("/account/activationcode/:uid", secure.ActivationCodeForm)
  router.Router.POST   ("/account/activationcode",      secure.ActivationCode)
  
  router.Router.GET    ("/account/passwordcode",      secure.PasswordCodeForm)
  router.Router.GET    ("/account/passwordcode/",     secure.PasswordCodeForm)
  router.Router.GET    ("/account/passwordcode/:uid", secure.PasswordCodeForm)
  router.Router.POST   ("/account/passwordcode",      secure.PasswordCode)
  
  router.Router.GET    ("/account/password/:uid", secure.PasswordForm)
  router.Router.PUT    ("/account/password",      secure.ChangePassword)
  
  router.Router.Handler("GET", "/captcha/*filepath", captcha.Server)
  router.Router.ServeFiles("/static/*filepath", http.Dir("./static"))

  go func(){
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      http.Redirect(w, r, "https://" + HTTP_HOST + HTTPS_PORT + r.URL.Path, http.StatusMovedPermanently)
    })
    log.Fatal(http.ListenAndServe(HTTP_HOST + HTTP_PORT,
      handlers.CombinedLoggingHandler(os.Stdout,
    http.DefaultServeMux)))
  }()

  log.Fatal(http.ListenAndServeTLS(HTTP_HOST + HTTPS_PORT, "cert.pem", "key.pem",
    context.ClearHandler(
    secure.AuthenticationHandler(
    handlers.HTTPMethodOverrideHandler(
    handlers.CompressHandler(
    handlers.CombinedLoggingHandler(os.Stdout, 
  router.Router)))))))
}
