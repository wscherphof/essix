package main

import (
  "log"
  "net/http"
  "html/template"
  "github.com/yossi/ace"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/msg"
  "bytes"
  "io"
)

var aceOptions = &ace.Options{ // var, no const, since the compiler says this literal isn't a constant
  BaseDir: "templates",
  FuncMap: template.FuncMap{
    "Msg": msg.Msg,
  },
}

func t (base string, inner string, data map[string]interface{}) (func(io.Writer, *http.Request)) {
  if data == nil {
    data = map[string]interface{}{}
  }
  return func(w io.Writer, r *http.Request) {
    language := msg.Language(r.Header.Get("Accept-Language"))
    data["lang"] = language
    if inner == "lang" {
      inner = base + "-" + language.Main
    }
    if tpl, err := ace.Load(base, inner, aceOptions); err != nil {
      log.Panicln("ERROR: ace.Load:", err)
    } else if err := tpl.Execute(w, data); err != nil {
      log.Panicln("ERROR: tpl.Execute:", err)
    }
  }
}

func T (base string, inner string, data map[string]interface{}) httprouter.Handle {
  return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    t(base, inner, data)(w, r)
  }
}

func TS (base string, inner string, data map[string]interface{}) (func(*http.Request)(string)) {
  var b *bytes.Buffer
  return func(r *http.Request) (string) {
    t(base, inner, data)(b, r)
    return string(b.Bytes())
  }
}
