package main

import (
  "net/http"
  "log"
  "os"
  "github.com/julienschmidt/httprouter"
  "github.com/gorilla/handlers"
  "github.com/gorilla/context"
)

func main () {
  DefineMessages()
  InitSecure()
  router := httprouter.New()

  router.GET("/", T("home", "", map[string]string{
  }))
  
  // TODO: https
  
  router.GET("/session", LogInForm)
  router.POST("/session", LogIn)
  router.DELETE("/session", LogOut)

  router.GET("/protected", Protected)
  
  router.ServeFiles("/static/*filepath", http.Dir("./static"))

  log.Fatal(http.ListenAndServe(":9090", 
  context.ClearHandler(
  handlers.HTTPMethodOverrideHandler(
  handlers.CombinedLoggingHandler(os.Stdout, 
  router)))))
}
