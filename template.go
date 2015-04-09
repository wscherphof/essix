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
    language := msg.Language(r.Header.Get("Accept-Language"))
    if inner == "lang" {
      inner = base + "-" + language.Main
    }
    // TODO: make the options a const
    tpl, err := ace.Load(base, inner, &ace.Options{
      BaseDir: "templates",
      // TODO: maybe move the Msg function to data to get it to enclose the language parameter again
      FuncMap: template.FuncMap{
        "Msg": msg.Msg,
      },
    })
    if err != nil {
      fmt.Println("ace.Load: ", err)
      return
    }
    if data == nil {
      data = map[string]interface{}{}
    }
    data["lang"] = language
    err = tpl.Execute(w, data)
    if err != nil {
      fmt.Println("tpl.Execute: ", err)
      return
    }
  }
}
