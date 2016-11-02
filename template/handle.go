package template

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Handle returns a Handle executing a template
func Handle(dir, base, inner string, opt_data ...map[string]interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var data map[string]interface{}
		if len(opt_data) == 1 {
			data = opt_data[0]
		}
		response(w, r, dir, base, inner, data)
	}
}
