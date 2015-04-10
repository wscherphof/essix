package main

import (
  "net/http"
  "log"
  "os"
  "github.com/julienschmidt/httprouter"
  "github.com/gorilla/handlers"
  "github.com/gorilla/context"
  "github.com/wscherphof/expeertise/db"
  "github.com/wscherphof/expeertise/secure"
  "github.com/wscherphof/expeertise/model"
)

func main () {
  db.Init("localhost:28015", "expeertise")
  secure.Init()
  model.Init()
  DefineMessages()
  router := httprouter.New()

  router.GET("/", T("home", "", nil))
  
  // TODO: https
  
  router.GET("/session", LogInForm)
  router.POST("/session", LogIn)
  router.DELETE("/session", LogOut)
  
  router.GET("/account", SignUpForm)
  router.POST("/account", SignUp)

  router.GET("/protected", Protected)
  
  router.ServeFiles("/static/*filepath", http.Dir("./static"))

  log.Fatal(http.ListenAndServe(":9090", 
  context.ClearHandler(
  handlers.HTTPMethodOverrideHandler(
  handlers.CombinedLoggingHandler(os.Stdout, 
  router)))))
}
