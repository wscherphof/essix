package router

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var (
	Router                 = httprouter.New()
	ErrInternalServerError = errors.New("ErrInternalServerError")
)

// ErrorHandle is a function that can be registered to a route to handle HTTP
// requests. Like httprouter.Handle, but returns an error value.
// Use this package's GET, PUT, POST, DELETE, PATH, OPTIONS, or HEAD to either
// execute a template, or return an Error. If Error isn't nil, an error template
// is executed.
type ErrorHandle func(http.ResponseWriter, *http.Request, httprouter.Params) *Error

func handleError(errorHandle ErrorHandle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if err := errorHandle(w, r, ps); err != nil {
			code := http.StatusInternalServerError
			if err.Conflict {
				code = http.StatusConflict
			} else {
				log.Printf("ERROR: %+v: %s %#v", r.URL, err.Error, err.Error)
				err.Error = ErrInternalServerError
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
				inner = "../" + err.Tail.dir + "/" + err.Tail.name
			}
			// Set the Content-Type to prevent CompressHandler from doing so after our WriteHeader()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(code)
			Template("router", "error", inner, data)(w, r, ps)
		}
	}
}

// GET registers a new request handle with the given path method GET.
func GET(path string, errorHandle ErrorHandle) { Router.Handle("GET", path, handleError(errorHandle)) }

// PUT registers a new request handle with the given path method PUT.
func PUT(path string, errorHandle ErrorHandle) { Router.Handle("PUT", path, handleError(errorHandle)) }

// POST registers a new request handle with the given path method POST.
func POST(path string, errorHandle ErrorHandle) { Router.Handle("POST", path, handleError(errorHandle)) }

// DELETE registers a new request handle with the given path method DELETE.
func DELETE(path string, errorHandle ErrorHandle) { Router.Handle("DELETE", path, handleError(errorHandle)) }

// PATCH registers a new request handle with the given path method PATCH.
func PATCH(path string, errorHandle ErrorHandle) { Router.Handle("PATCH", path, handleError(errorHandle)) }

// OPTIONS registers a new request handle with the given path method OPTIONS.
func OPTIONS(path string, errorHandle ErrorHandle) { Router.Handle("OPTIONS", path, handleError(errorHandle)) }

// HEAD registers a new request handle with the given path method HEAD.
func HEAD(path string, errorHandle ErrorHandle) { Router.Handle("HEAD", path, handleError(errorHandle)) }
