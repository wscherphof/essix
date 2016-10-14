/*
Package router provides the httprouter that registers the application's request
routes.
*/
package router

import (
	"github.com/julienschmidt/httprouter"
)

/*
Router is the httprouter for the application.
*/
var Router = httprouter.New()

// GET registers a handler for a GET request to the given path.
func GET(path string, handle httprouter.Handle)     { Router.GET(path, handle) }

// PUT registers a handler for a PUT request to the given path.
func PUT(path string, handle httprouter.Handle)     { Router.PUT(path, handle) }

// POST registers a handler for a POST request to the given path.
func POST(path string, handle httprouter.Handle)    { Router.POST(path, handle) }

// DELETE registers a handler for a DELETE request to the given path.
func DELETE(path string, handle httprouter.Handle)  { Router.DELETE(path, handle) }

// HEAD registers a handler for a HEAD request to the given path.
func HEAD(path string, handle httprouter.Handle)    { Router.HEAD(path, handle) }

// OPTIONS registers a handler for a OPTIONS request to the given path.
func OPTIONS(path string, handle httprouter.Handle) { Router.OPTIONS(path, handle) }

// PATCH registers a handler for a PATCH request to the given path.
func PATCH(path string, handle httprouter.Handle)   { Router.PATCH(path, handle) }
