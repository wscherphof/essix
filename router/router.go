package router

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var Router = httprouter.New()

func GET(path string, handle httprouter.Handle)     { Router.GET(path, handle) }
func PUT(path string, handle httprouter.Handle)     { Router.PUT(path, handle) }
func POST(path string, handle httprouter.Handle)    { Router.POST(path, handle) }
func DELETE(path string, handle httprouter.Handle)  { Router.DELETE(path, handle) }
func HEAD(path string, handle httprouter.Handle)    { Router.HEAD(path, handle) }
func OPTIONS(path string, handle httprouter.Handle) { Router.OPTIONS(path, handle) }
func PATCH(path string, handle httprouter.Handle)   { Router.PATCH(path, handle) }

func Redirect(w http.ResponseWriter, r *http.Request, data map[string]string) {
	var query string
	if data != nil {
		for key, value := range data {
		    query += key + "=" + value + "&"
		}		
	}
	var path = r.URL.Path + "/prg"
	if len(query) > 0 {
		path += "?" + query
	}
	http.Redirect(w, r, path, http.StatusSeeOther)
}
