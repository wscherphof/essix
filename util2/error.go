package util2

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "log"
)

type Error struct{
  Error error
  Conflict bool
  Tail string
  Data map[string]interface{}
}

func NewError (e error, tail ...string) (err *Error) {
  err = &Error{Error: e}
  if len(tail) > 0 {
    err.Tail = tail[0]
  }
  return
}

type ErrorHandle func(http.ResponseWriter, *http.Request, httprouter.Params)(*Error)

func ErrorHandleFunc (f ErrorHandle) (handle httprouter.Handle) {
  return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    if err := f(w, r, ps); err != nil {
      code := http.StatusInternalServerError
      if err.Conflict {
        code = http.StatusConflict
      }
      // Set the Content-Type to prevent CompressHandler from doing so after our WriteHeader()
      w.Header().Set("Content-Type", "text/html; charset=utf-8")
      w.WriteHeader(code)
      Template("error", "", map[string]interface{}{
        "error": err.Error.Error(),
      })(w, r, ps)
      if len(err.Tail) > 0 {
        Template(err.Tail + "_error-tail", "", err.Data)(w, r, ps)
      }
      if code >= 500 {
        log.Println("ERROR:", err.Error, "- Path:", r.URL.Path)
      }
    }
  }
}
