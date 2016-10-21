package example

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/secure"
	"github.com/wscherphof/essix/template"
	"<model>"
	"net/http"
	"strings"
)

func ProfileForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := secure.Authentication(r)
	t := template.GET(w, r, "example", "ProfileForm")
	if profile := readProfile(w, r, account.ID); profile != nil {
		t.Set("email", account.Email)
		t.Set("profile", profile)
		t.Set("countries", Countries())
		t.Run()
	}
}

func Profile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := secure.Authentication(r)
	if profile := readProfile(w, r, account.ID); profile != nil {
		profile.Country = r.FormValue("country")
		profile.Postcode = strings.ToUpper(r.FormValue("postcode"))
		profile.FirstName = r.FormValue("firstname")
		profile.LastName = r.FormValue("lastname")
		if err := profile.Update(profile); err != nil {
			template.Error(w, r, err, false)
		} else {
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		}
	}
}

func readProfile(w http.ResponseWriter, r *http.Request, id string) (profile *model.Profile) {
	profile = model.InitProfile(id)
	if err, empty := profile.Read(profile); err != nil && !empty {
		template.Error(w, r, err, false)
		profile = nil
	}
	return
}
