package util

import (
  "log"
  "net/http"
  "github.com/julienschmidt/httprouter"
)

func Error (w http.ResponseWriter, r *http.Request, ps httprouter.Params, err error, codes ...int) {
  if err == nil {
    return
  }
  code := http.StatusInternalServerError
  if len(codes) > 0 {
    code = codes[0]
  }
  render := BTemplate("error", "", map[string]interface{}{
    "error": err.Error(),
  })(r)
  w.WriteHeader(code)
  w.Write(render)
  if code >= 500 {
    log.Println("ERROR:", err, "- Path:", r.URL.Path)
  }
}
