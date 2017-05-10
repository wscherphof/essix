package template

import (
	"errors"
	"log"
	"net/http"
)

var ErrInternalServerError = errors.New("ErrInternalServerError")

/*
Error renders a generic template with details about the given error.

If conflict is true, the status code is http.StatusConflict. Otherwise, it's
http.StatusInternalServerError.
*/
func Error(w http.ResponseWriter, r *http.Request, err error, conflict bool) {
	errorTemplate(w, r, err, conflict, nil)
}

/*
ErrorTail is like Error, but tails the generic error template with another, as an
"inner" template.

The tail template should start with:
	= content tail

Any data is passed to the tail template.
*/
func ErrorTail(w http.ResponseWriter, r *http.Request, err error, conflict bool, dir, tail string, opt_data ...map[string]interface{}) {
	var data map[string]interface{}
	if len(opt_data) == 1 {
		data = opt_data[0]
	}
	errorTemplate(w, r, err, conflict, data, dir, tail)
}

func errorTemplate(w http.ResponseWriter, r *http.Request, err error, conflict bool, errData map[string]interface{}, tail ...string) {
	status := http.StatusInternalServerError
	if conflict {
		status = http.StatusConflict
	} else {
		log.Printf("ERROR: %s %+v: %s %#v", r.Method, r.URL, err, err)
		err = ErrInternalServerError
	}
	data := make(map[string]interface{}, 2+len(errData))
	data["error"] = err.Error()
	data["url"] = r.URL.String()
	if errData != nil {
		for k, v := range errData {
			data[k] = v
		}
	}
	inner := ""
	if len(tail) == 2 {
		inner = "../" + tail[0] + "/" + tail[1]
	}
	response(w, r, "template", "Error", inner, data, status)
}
