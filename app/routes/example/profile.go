package example

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/secure"
	"github.com/wscherphof/essix/template"
	"<model>"
	"net/http"
	"time"
	"log"
)

func ProfileForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := secure.Authentication(r)
	t := template.GET(w, r, "example", "ProfileForm")
	if profile := readProfile(w, r, account.ID); profile != nil {
		t.Set("email", account.Email)
		t.Set("profile", profile)
		t.Set("countries", Countries())
		t.Set("timezones", TimeZones())
		if profile.TimeZone != "" {
			if location, err := time.LoadLocation(profile.TimeZone); err != nil {
				log.Println("WARNING:", err)
			} else {
				profile.Modified = profile.Modified.In(location)
			}
		}
		t.Run()
	}
}

func Profile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := secure.Authentication(r)
	if profile := readProfile(w, r, account.ID); profile != nil {
		profile.Country = r.FormValue("country")
		profile.TimeZone = r.FormValue("timezone")
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
