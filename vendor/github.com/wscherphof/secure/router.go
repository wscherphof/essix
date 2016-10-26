package secure

import (
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/url"
)

var router *SecureRouter

/*
Router returns the secure router.
*/
func Router() *SecureRouter {
	if router == nil {
		router = &SecureRouter{httprouter.New()}
	}
	return router
}

/*
A SecureRouter is a secured httprouter.
PUT, POST, PATCH, and DELETE handles check for a valid FormToken encrypted token
string in the request's "_formtoken" FormValue.
*/
type SecureRouter struct {
	*httprouter.Router
}

// PUT registers a handler for a PUT request to the given path.
// The handler is only run if the request carries a valid form token.
func (r *SecureRouter) PUT(path string, handle httprouter.Handle) {
	r.Handle("PUT", path, formTokenHandle(handle))
}

// POST registers a handler for a POST request to the given path.
// The handler is only run if the request carries a valid form token.
func (r *SecureRouter) POST(path string, handle httprouter.Handle) {
	r.Handle("POST", path, formTokenHandle(handle))
}

// PATCH registers a handler for a PATCH request to the given path.
// The handler is only run if the request carries a valid form token.
func (r *SecureRouter) PATCH(path string, handle httprouter.Handle) {
	r.Handle("PATCH", path, formTokenHandle(handle))
}

// DELETE registers a handler for a DELETE request to the given path.
// The handler is only run if the request carries a valid form token.
func (r *SecureRouter) DELETE(path string, handle httprouter.Handle) {
	r.Handle("DELETE", path, formTokenHandle(handle))
}

// GET registers a handler for a GET request to the given path.
func (r *SecureRouter) GET(path string, handle httprouter.Handle) {
	r.Handle("GET", path, clearHandle(handle))
}

// HEAD registers a handler for a HEAD request to the given path.
func (r *SecureRouter) HEAD(path string, handle httprouter.Handle) {
	r.Handle("HEAD", path, clearHandle(handle))
}

// OPTIONS registers a handler for a OPTIONS request to the given path.
func (r *SecureRouter) OPTIONS(path string, handle httprouter.Handle) {
	r.Handle("OPTIONS", path, clearHandle(handle))
}

func clearHandle(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handle(w, r, ps)
		context.Clear(r)
	}
}

/*
FormValueName is the name of the FormValue that is checked in PUT, POST, PATCH,
and DELETE handles.
*/
const FormValueName = "_formtoken"

func formTokenHandle(handle httprouter.Handle) httprouter.Handle {
	return clearHandle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		this, that := NewFormToken(r), new(FormToken)
		if err := that.Parse(r.FormValue(FormValueName)); err != nil {
			log.Printf("WARNING: %s %s %s", err, this.IP, this.Path)
		} else {
			referer, _ := url.Parse(r.Referer())
			// Timestamp not considered, since key rotation will outdate old tokens automatically
			if that.IP != this.IP || (that.Path != this.Path && that.Path != referer.Path) {
				log.Printf("WARNING: Form token invalid %s %s", this.IP, this.Path)
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`<!DOCTYPE html>
					<html>
						<head>
							<meta charset="utf-8">
						</head>
						<body>
							<h2>Form token validation failed</h2>
							<a id="location" href="` + referer.Path + `">Back</a>
						</body>
					</html>
				`))
			} else {
				handle(w, r, ps)
			}
		}
	})
}

/*
secure.Handle ensures the client is logged in when accessing a certian route,
redirecting to the log in page if not. The given Handle function should call
Authentication() to get the client's account details.

If the cookie is missing, the session has timed out, or the cookie data is
invalidated though the ValidateCookie function, the response then gets status
403 Forbidden, and the browser will redirect to config.LogInPath.
*/
func Handle(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if authenticate(w, r, false) {
			handle(w, r, ps)
		}
	}
}

/*
secure.IfHandle calls the one Hanlde function for logged-in clients, and the
other for logged-out clients.
*/
func IfHandle(authenticatedHandle httprouter.Handle, unauthenticatedHandle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if authenticate(w, r, true) {
			authenticatedHandle(w, r, ps)
		} else {
			unauthenticatedHandle(w, r, ps)
		}
	}
}
