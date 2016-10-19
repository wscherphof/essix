/*
Package router provides the httprouter that registers the application's request
routes.
*/
package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/secure"
	"net/http"
)

/*
Router is the httprouter for the application.
*/
var Router = httprouter.New()

func formTokenHandle(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if err := secure.ValidateFormToken(r); err == nil {
			handle(w, r, ps)
		}
	}
}

// GET registers a handler for a GET request to the given path.
func GET(path string, handle httprouter.Handle) { Router.GET(path, handle) }

// PUT registers a handler for a PUT request to the given path.
// The handler is only run if the request carries a valid form token.
func PUT(path string, handle httprouter.Handle) { Router.PUT(path, formTokenHandle(handle)) }

// POST registers a handler for a POST request to the given path.
// The handler is only run if the request carries a valid form token.
func POST(path string, handle httprouter.Handle) { Router.POST(path, formTokenHandle(handle)) }

// DELETE registers a handler for a DELETE request to the given path.
// The handler is only run if the request carries a valid form token.
func DELETE(path string, handle httprouter.Handle) { Router.DELETE(path, formTokenHandle(handle)) }

// HEAD registers a handler for a HEAD request to the given path.
func HEAD(path string, handle httprouter.Handle) { Router.HEAD(path, handle) }

// OPTIONS registers a handler for a OPTIONS request to the given path.
func OPTIONS(path string, handle httprouter.Handle) { Router.OPTIONS(path, handle) }

// PATCH registers a handler for a PATCH request to the given path.
// The handler is only run if the request carries a valid form token.
func PATCH(path string, handle httprouter.Handle) { Router.PATCH(path, formTokenHandle(handle)) }
