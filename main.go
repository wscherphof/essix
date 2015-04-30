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
  "github.com/wscherphof/expeertise/util"
)

func main () {
  db.Init("localhost:28015", "expeertise")
  config.Init()
  secure.Init()
  model.Init()
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
  
  router.ServeFiles("/static/*filepath", http.Dir("./static"))

  log.Fatal(http.ListenAndServe(":9090", 
  context.ClearHandler(
  handlers.HTTPMethodOverrideHandler(
  handlers.CombinedLoggingHandler(os.Stdout, 
  router)))))
}
