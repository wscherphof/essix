package main

import (
  "fmt"
  "net/http"
  "html/template"
  "github.com/yossi/ace"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/msg"
)

func T (base string, inner string, data map[string]string) func (http.ResponseWriter, *http.Request, httprouter.Params) {
  return func (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    Msg, language := msg.Language(r.Header.Get("Accept-Language"))
    data["lang"] = language
    tpl, err := ace.Load(base, inner, &ace.Options{
      BaseDir: "templates",
      FuncMap: template.FuncMap{
        "Msg": Msg,
      },
    })
    if err != nil {
        fmt.Println(err)
        return
    }
    err = tpl.Execute(w, data)
    if err != nil {
        fmt.Println(err)
        return
    }
  }
}
