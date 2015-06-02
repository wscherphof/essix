package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/expeertise/util"
	"log"
	"net/http"
)

var Router = httprouter.New()

type Error struct {
	Error    error
	Conflict bool
	Tail     string
	Data     map[string]interface{}
}

func NewError(e error, tail ...string) (err *Error) {
	err = &Error{Error: e}
	if len(tail) > 0 {
		err.Tail = tail[0]
	}
	return
}

type ErrorHandle func(http.ResponseWriter, *http.Request, httprouter.Params) *Error

func ErrorHandleFunc(errorHandle ErrorHandle) (handle httprouter.Handle) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if err := errorHandle(w, r, ps); err != nil {
			code := http.StatusInternalServerError
			if err.Conflict {
				code = http.StatusConflict
			}
			// Set the Content-Type to prevent CompressHandler from doing so after our WriteHeader()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(code)
			Template("error", "", map[string]interface{}{
				"Error": err.Error,
			})(w, r, ps)
			if len(err.Tail) > 0 {
				Template(err.Tail+"_error-tail", "", err.Data)(w, r, ps)
			}
			if code >= 500 {
				log.Println("ERROR:", err.Error, "- Path:", r.URL.Path)
			}
		}
	}
}

func Handle(method, path string, handle ErrorHandle) {
	Router.Handle(method, path, ErrorHandleFunc(handle))
}
func GET(path string, handle ErrorHandle)     { Handle("GET", path, handle) }
func PUT(path string, handle ErrorHandle)     { Handle("PUT", path, handle) }
func POST(path string, handle ErrorHandle)    { Handle("POST", path, handle) }
func DELETE(path string, handle ErrorHandle)  { Handle("DELETE", path, handle) }
func PATCH(path string, handle ErrorHandle)   { Handle("PATCH", path, handle) }
func OPTIONS(path string, handle ErrorHandle) { Handle("OPTIONS", path, handle) }
func HEAD(path string, handle ErrorHandle)    { Handle("HEAD", path, handle) }

func Template(base string, inner string, data map[string]interface{}) ErrorHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *Error) {
		// TODO: try if we can do ps.ByName() from the ace template..
		util.Template(base, inner, data)(w, r)
		return
	}
}
