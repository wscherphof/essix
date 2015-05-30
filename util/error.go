package util

import (
  "log"
  "net/http"
  "github.com/julienschmidt/httprouter"
)

// TODO: rename to Catch(Error)
func Handle (w http.ResponseWriter, r *http.Request, ps httprouter.Params) func(error, bool, string, map[string]interface{}) {
  return func(err error, conflict bool, tail string, data map[string]interface{}) {
    if err == nil {
      return
    }
    code := http.StatusInternalServerError
    if conflict {
      code = http.StatusConflict
    }
    // Set the Content-Type to prevent CompressHandler from doing so after our WriteHeader()
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(code)
    Template("error", "", map[string]interface{}{
      "error": err.Error(),
    })(w, r, ps)
    if len(tail) > 0 {
      Template(tail + "_error-tail", "", data)(w, r, ps)
    }
    if code >= 500 {
      log.Println("ERROR:", err, "- Path:", r.URL.Path)
    }
  }
}
