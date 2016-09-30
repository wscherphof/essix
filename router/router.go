package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/util"
	"log"
	"net/http"
	"errors"
)

var (
	Router                 = httprouter.New()
	ErrInternalServerError = errors.New("ErrInternalServerError")
)

// TemplateHandle returns a Handle executing a template
func TemplateHandle(dir, base, inner string, data map[string]interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		util.Template(w, r, dir, base, inner, data)
	}
}

// Error executes a template reporting on e
func Error(w http.ResponseWriter, r *http.Request, err error, conflict bool, tail ...string) {
	errorTemplate(w, r, err, conflict, nil, tail...)
}

// DataError executes a template reporting on e & data
func DataError(w http.ResponseWriter, r *http.Request, err error, data map[string]interface{}, tail ...string) {
	errorTemplate(w, r, err, true, data, tail...)
}

func errorTemplate(w http.ResponseWriter, r *http.Request, err error, conflict bool, errData map[string]interface{}, tail ...string) {
	code := http.StatusInternalServerError
	if conflict {
		code = http.StatusConflict
	} else {
		log.Printf("ERROR: %+v: %s %#v", r.URL, err, err)
		err = ErrInternalServerError
	}
	data := map[string]interface{}{
		"Error": err,
		"Path":  r.URL.Path,
	}
	if errData != nil {
		for k, v := range errData {
			data[k] = v
		}
	}
	inner := ""
	if len(tail) == 2 {
		inner = "../" + tail[0] + "/" + tail[1] + "_error-tail"
	}
	// Set the Content-Type to prevent CompressHandler from doing so after our WriteHeader()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	util.Template(w, r, "router", "error", inner, data)
}
