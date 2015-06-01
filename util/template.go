package util

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

func aceOptions (r *http.Request) (*ace.Options) {
  return &ace.Options{
    BaseDir: "templates",
    FuncMap: template.FuncMap{
      "Msg": msg.Msg(r),
    },
  }
}

func t (base string, inner string, data map[string]interface{}) (func(io.Writer, *http.Request)) {
  if data == nil {
    data = map[string]interface{}{}
  }
  return func(w io.Writer, r *http.Request) {
    lang := msg.Language(r)
    data["lang"] = lang
    if inner == "lang" {
      inner = base + "-" + lang.Main
    }
    if tpl, err := ace.Load(base, inner, aceOptions(r)); err != nil {
      log.Panicln("ERROR: ace.Load:", err)
    } else if err := tpl.Execute(w, data); err != nil {
      log.Panicln("ERROR: tpl.Execute:", err)
    }
  }
} 

func Template (base string, inner string, data map[string]interface{}) func(http.ResponseWriter, *http.Request, httprouter.Params)(*Error) {
  return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *Error) {
    // TODO: try if we can do ps.ByName() from the ace template..
    t(base, inner, data)(w, r)
    return
  }
}

func BTemplate (base string, inner string, data map[string]interface{}) (func(*http.Request)([]byte)) {
  var b bytes.Buffer
  return func(r *http.Request) ([]byte) {
    t(base, inner, data)(&b, r)
    return b.Bytes()
  }
}
