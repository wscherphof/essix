package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/util"
	"net/http"
)

// Template
func Template(dir, base, inner string, data map[string]interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		util.Template(dir, base, inner, data)(w, r)
	}
}
