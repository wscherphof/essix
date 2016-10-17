package example

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/secure"
	"github.com/wscherphof/essix/template"
	"net/http"
	"strings"
)

func ProfileForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := secure.Authentication(w, r)
	template.Run(w, r, "example", "profile", "", map[string]interface{}{
		"Account":   acc,
		"Countries": data.Countries(),
		"Initial":   (acc.ValidateFields() != nil),
	})
}

func Profile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	acc := secure.Authentication(w, r)
	// TODO: fetch profile
	// TODO: update profile
	// acc.Country = r.FormValue("country")
	// acc.Postcode = strings.ToUpper(r.FormValue("postcode"))
	// acc.FirstName = r.FormValue("firstname")
	// acc.LastName = r.FormValue("lastname")
	// POST, redirect, GET
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}
