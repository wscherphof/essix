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
  "github.com/wscherphof/expeertise/util2"
)

const (
  HTTP_HOST  = "localhost"
  HTTP_PORT  = ":9090"
  HTTPS_PORT = ":10443"
  DB_HOST    = "localhost"
  DB_PORT    = ":28015"
  DB_NAME    = "expeertise"
)

var router = httprouter.New()

func errorHandle (method, pattern string, handle util2.ErrorHandle) {
  router.Handle(method, pattern, util2.ErrorHandleFunc(handle))
}
func GET    (pattern string, handle util2.ErrorHandle) {errorHandle("GET",    pattern, handle)}
func PUT    (pattern string, handle util2.ErrorHandle) {errorHandle("PUT",    pattern, handle)}
func POST   (pattern string, handle util2.ErrorHandle) {errorHandle("POST",   pattern, handle)}
func DELETE (pattern string, handle util2.ErrorHandle) {errorHandle("DELETE", pattern, handle)}

func main () {
  db.Init(DB_HOST + DB_PORT, DB_NAME)
  config.Init()
  secure.Init()
  model.Init()
  captcha.Init()
  DefineMessages()

  // TODO: differentiate whether logged in
  GET    ("/", util2.Template("home", "", nil))
  
  // TODO: sign up w/ just email & pwd; then on first login, ask further details
  // TODO: change email address (only when logged in, but still w/ an email to the new address)
  router.GET    ("/account", secure.AccountForm)
  POST   ("/account", secure.SignUp)
  PUT    ("/account", secure.SecureHandle(secure.UpdateAccount))
  // TODO: router.DELETE ("/account", secure.Authenticate(secure.TerminateAccount))

  router.GET    ("/session", secure.LogInForm)
  router.POST   ("/session", secure.LogIn)
  router.DELETE ("/session", secure.LogOut)

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
    context.ClearHandler(
    secure.AuthenticationHandler(
    handlers.HTTPMethodOverrideHandler(
    handlers.CompressHandler(
    handlers.CombinedLoggingHandler(os.Stdout, 
  router)))))))
}
