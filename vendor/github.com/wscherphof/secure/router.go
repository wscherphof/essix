package secure

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/url"
)

/*
Router is a new secured httprouter.
PUT, POST, PATCH, and DELETE handles check for a valid FormToken encrypted token
string in the request's "_formtoken" FormValue.
*/
var Router = &router{httprouter.New()}

type router struct {
	*httprouter.Router
}

// PUT registers a handler for a PUT request to the given path.
// The handler is only run if the request carries a valid form token.
func (r *router) PUT(path string, handle httprouter.Handle) {
	r.Handle("PUT", path, formTokenHandle(handle))
}

// POST registers a handler for a POST request to the given path.
// The handler is only run if the request carries a valid form token.
func (r *router) POST(path string, handle httprouter.Handle) {
	r.Handle("POST", path, formTokenHandle(handle))
}

// PATCH registers a handler for a PATCH request to the given path.
// The handler is only run if the request carries a valid form token.
func (r *router) PATCH(path string, handle httprouter.Handle) {
	r.Handle("PATCH", path, formTokenHandle(handle))
}

// DELETE registers a handler for a DELETE request to the given path.
// The handler is only run if the request carries a valid form token.
func (r *router) DELETE(path string, handle httprouter.Handle) {
	r.Handle("DELETE", path, formTokenHandle(handle))
}

/*
FormValueName is the name of the FormValue that is checked in PUT, POST, PATCH,
and DELETE handles.
*/
const FormValueName = "_formtoken"

func formTokenHandle(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		this, that := NewFormToken(r), new(FormToken)
		if err := that.Parse(r.FormValue(FormValueName)); err != nil {
			log.Printf("WARNING: %s %s %s", err, this.IP, this.Path)
		} else {
			referer, _ := url.Parse(r.Referer())
			// Timestamp not considered, since key rotation will outdate old tokens automatically
			if that.IP != this.IP || (that.Path != this.Path && that.Path != referer.Path) {
				log.Printf("WARNING: Form token invalid %s %s", this.IP, this.Path)
			} else {
				handle(w, r, ps)
			}
		}
	}
}

/*
secure.Handle ensures the client is logged in when accessing a certian route,
redirecting to the log in page if not. The given Handle function should call
Authentication() to get the client's account details.
*/
func Handle(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if Authentication(w, r, false) != nil {
			handle(w, r, ps)
		}
	}
}

/*
secure.IfHandle calls the one Hanlde function for logged-in clients, and the
other for logged-out clients.
*/
func IfHandle(authenticated httprouter.Handle, unauthenticated httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if Authentication(w, r, true) != nil {
			authenticated(w, r, ps)
		} else {
			unauthenticated(w, r, ps)
		}
	}
}
