package router

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/expeertise/util"
	"log"
	"net/http"
)

var (
	Router                 = httprouter.New()
	ErrInternalServerError = errors.New("Internal server error")
)

type tailType struct {
	dir  string
	name string
}

type Error struct {
	Error    error
	Conflict bool
	Tail     *tailType
	Data     map[string]interface{}
}

type ErrorHandle func(http.ResponseWriter, *http.Request, httprouter.Params) *Error

func ErrorHandleFunc(errorHandle ErrorHandle) (handle httprouter.Handle) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if err := errorHandle(w, r, ps); err != nil {
			code := http.StatusInternalServerError
			if err.Conflict {
				code = http.StatusConflict
			} else {
				log.Printf("ERROR: %+v: %s %#v", r.URL, err.Error, err.Error)
				err.Error = ErrInternalServerError
			}
			// Set the Content-Type to prevent CompressHandler from doing so after our WriteHeader()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(code)
			Template("router", "error", "", map[string]interface{}{
				"Error": err.Error,
			})(w, r, ps)
			if err.Tail != nil {
				Template(err.Tail.dir, err.Tail.name, "", err.Data)(w, r, ps)
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

func Template(dir, base, inner string, data map[string]interface{}) ErrorHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *Error) {
		// TODO: try if we can do ps.ByName() from the ace template..
		util.Template(dir, base, inner, data)(w, r)
		return
	}
}

func NewError(e error, tail ...string) (err *Error) {
	if e != nil {
		err = &Error{Error: e}
		if len(tail) == 2 {
			err.Tail = &tailType{
				dir:  tail[0],
				name: tail[1] + "_error-tail",
			}
		}
	}
	return
}

func IfError(e error, tail ...string) (err *Error) {
	return NewError(e, tail...)
}
