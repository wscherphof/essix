package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/util"
	"net/http"
)

var	Router                 = httprouter.New()

// TemplateHandle returns a Handle executing a template
func TemplateHandle(dir, base, inner string, data map[string]interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		util.Template(w, r, dir, base, inner, data)
	}
}
