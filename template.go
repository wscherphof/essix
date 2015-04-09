package main

import (
  "fmt"
  "net/http"
  "html/template"
  "github.com/yossi/ace"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/msg"
)

func T (base string, inner string, data map[string]interface{}) func (http.ResponseWriter, *http.Request, httprouter.Params) {
  return func (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    if data == nil {
      data = map[string]interface{}{}
    }
    Msg, language := msg.Language(r.Header.Get("Accept-Language"))
    data["lang"] = language
    if inner == "lang" {
      inner = base + "-" + language.Main
    }
    tpl, err := ace.Load(base, inner, &ace.Options{
      BaseDir: "templates",
      // TODO: make lang a parameter to Msg, to prevent cache hit of other language
      FuncMap: template.FuncMap{
        "Msg": Msg,
      },
    })
    if err != nil {
        fmt.Println("ace.Load: ", err)
        return
    }
    err = tpl.Execute(w, data)
    if err != nil {
        fmt.Println("tpl.Execute: ", err)
        return
    }
  }
}
