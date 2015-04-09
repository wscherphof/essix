package main

import (
  "fmt"
  "net/http"
  "html/template"
  "github.com/yossi/ace"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/msg"
)

var aceOptions = &ace.Options{ // var, no const, since the compiler says this literal isn't a constant
  BaseDir: "templates",
  FuncMap: template.FuncMap{
    "Msg": msg.Msg,
  },
}

func T (base string, inner string, data map[string]interface{}) func (http.ResponseWriter, *http.Request, httprouter.Params) {
  if data == nil {
    data = map[string]interface{}{}
  }
  return func (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    language := msg.Language(r.Header.Get("Accept-Language"))
    if inner == "lang" {
      inner = base + "-" + language.Main
    }
    tpl, err := ace.Load(base, inner, aceOptions)
    if err != nil {
      fmt.Println("ace.Load: ", err)
      return
    }
    data["lang"] = language
    err = tpl.Execute(w, data)
    if err != nil {
      fmt.Println("tpl.Execute: ", err)
      return
    }
  }
}
