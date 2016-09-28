package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/util"
	"net/http"
)

func Template(dir, base, inner string, data map[string]interface{}) ErrorHandle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *Error) {
		util.Template(dir, base, inner, data)(w, r)
		return
	}
}
