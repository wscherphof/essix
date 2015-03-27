package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "net/http/httptest"
)

func Error (w http.ResponseWriter, r *http.Request, ps httprouter.Params, err error, codes ...int) {
  if err == nil {
    return
  }
  code := http.StatusInternalServerError
  if len(codes) > 0 {
    code = codes[0]
  }
  rec := httptest.NewRecorder()
  T("error", "", map[string]string{
    "error": err.Error(),
  })(rec, r, ps)
  w.WriteHeader(code)
  w.Write(rec.Body.Bytes())
}
