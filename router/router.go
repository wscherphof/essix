package router

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var (
	Router                 = httprouter.New()
	errInternalServerError = errors.New("errInternalServerError")
)

type ErrorHandle func(http.ResponseWriter, *http.Request, httprouter.Params) *Error

func errorHandleFunc(errorHandle ErrorHandle) (handle httprouter.Handle) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if err := errorHandle(w, r, ps); err != nil {
			code := http.StatusInternalServerError
			if err.Conflict {
				code = http.StatusConflict
			} else {
				log.Printf("ERROR: %+v: %s %#v", r.URL, err.Error, err.Error)
				err.Error = errInternalServerError
			}
			data := map[string]interface{}{
				"Error": err.Error,
				"Path":  r.URL.Path,
			}
			if err.Data != nil {
				for k, v := range err.Data {
					data[k] = v
				}
			}
			inner := ""
			if err.Tail != nil {
				inner = "../../" + err.Tail.dir + "/templates/" + err.Tail.name
			}
			// Set the Content-Type to prevent CompressHandler from doing so after our WriteHeader()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(code)
			Template("router", "error", inner, data)(w, r, ps)
		}
	}
}

func handle(method, path string, h ErrorHandle) {
	Router.Handle(method, path, errorHandleFunc(h))
}
func GET(path string, h ErrorHandle)     { handle("GET", path, h) }
func PUT(path string, h ErrorHandle)     { handle("PUT", path, h) }
func POST(path string, h ErrorHandle)    { handle("POST", path, h) }
func DELETE(path string, h ErrorHandle)  { handle("DELETE", path, h) }
func PATCH(path string, h ErrorHandle)   { handle("PATCH", path, h) }
func OPTIONS(path string, h ErrorHandle) { handle("OPTIONS", path, h) }
func HEAD(path string, h ErrorHandle)    { handle("HEAD", path, h) }
