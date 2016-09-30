package template

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Handle returns a Handle executing a template
func Handle(dir, base, inner string, data map[string]interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		Run(w, r, dir, base, inner, data)
	}
}
