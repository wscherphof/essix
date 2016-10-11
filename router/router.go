package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/template"
	"net/http"
	"strings"
)

var Router = httprouter.New()

func GET(path string, handle httprouter.Handle)     { Router.GET(path, handle) }
func PUT(path string, handle httprouter.Handle)     { Router.PUT(path, handle) }
func POST(path string, handle httprouter.Handle)    { Router.POST(path, handle) }
func DELETE(path string, handle httprouter.Handle)  { Router.DELETE(path, handle) }
func HEAD(path string, handle httprouter.Handle)    { Router.HEAD(path, handle) }
func OPTIONS(path string, handle httprouter.Handle) { Router.OPTIONS(path, handle) }
func PATCH(path string, handle httprouter.Handle)   { Router.PATCH(path, handle) }

func PRG(w http.ResponseWriter, r *http.Request, dir, base, inner string, keys ...string) (prg func(...string)) {
	switch r.Method {
	case "GET":
		data := make(map[string]interface{})
		for i := 0; i < len(keys); i++ {
			data[keys[i]] = r.FormValue(keys[i])
		}
		template.Run(w, r, dir, base, inner, data)
	case "PUT", "POST", "DELETE":
		prg = func(values ...string) {
			var query string
			for i := 0; i < len(values); i++ {
				query += keys[i] + "=" + values[i] + "&"
			}
			var path = r.URL.Path + "/" + strings.ToLower(r.Method)
			if len(query) > 0 {
				path += "?" + query
			}
			http.Redirect(w, r, path, http.StatusSeeOther)
		}
	}
	return
}
