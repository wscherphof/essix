package main

import (
  "net/http"
  "log"
  "os"
  "github.com/julienschmidt/httprouter"
  "github.com/gorilla/handlers"
)

func main () {
  DefineMessages()
  InitSecure()
  router := httprouter.New()

  router.GET("/", T("base", "", map[string]string{
    "Msg": "Hello Ace",
  }))
  
  // TODO: https
  
  router.GET("/login", LoginForm)
  router.POST("/login", Login)

  router.GET("/protected", Protected)
  
  router.ServeFiles("/static/*filepath", http.Dir("./static"))

  log.Fatal(http.ListenAndServe(":9090", handlers.CombinedLoggingHandler(os.Stdout, router)))
}
