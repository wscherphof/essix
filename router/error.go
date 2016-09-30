package router

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type tailType struct {
	dir  string
	name string
}

// error
type routerError struct {
	Error    error
	Conflict bool
	Tail     *tailType
	Data     map[string]interface{}
}

// newError constructs a routerError
func newError(e error, conflict bool, tail ...string) (err *routerError) {
	err = &routerError{
		Error: e,
		Conflict: conflict,
	}
	if len(tail) == 2 {
		err.Tail = &tailType{
			dir:  tail[0],
			name: tail[1] + "_error-tail",
		}
	}
	return
}

// Error returns a Handle template reporting on e
func Error(e error, conflict bool, tail ...string) httprouter.Handle {
	err := newError(e, conflict, tail...)
	return errorHandle(err)
}

// DataError returns a Handle template reporting on e & data
func DataError(e error, data map[string]interface{}, tail ...string) httprouter.Handle {
	err := newError(e, true, tail...)
	err.Data = data
	return errorHandle(err)
}

func errorHandle(err *routerError) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
