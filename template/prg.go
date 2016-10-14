package template

import (
	"net/http"
	"net/url"
	"strings"
)

type baseType struct {
	*url.Values
	w http.ResponseWriter
	r *http.Request
}

type PRGType struct {
	*baseType
}

// Run redirects to path/method, including set values.
func (t *PRGType) Run() {
	path := t.r.URL.Path
	path += "/" + strings.ToLower(t.r.Method)
	path += "?" + t.Values.Encode()
	http.Redirect(t.w, t.r, path, http.StatusSeeOther)
}

/*
PRG implments the Post-Redirect-Get pattern.

If the request's method is GET, it returns nil, and runs the template,
populating it with the data stored in the POST processing.

Otherwise, it returns a type that, just as net/url.Values, listens to Set(key, value)
to register string data. Call Run() to execute the redirect. Register a GET route
for the redirect with the same path, suffixed with /<method>

Example:
	router.POST("/account", ratelimit.Handle(account.NewAccount))
	router.GET("/account/post", account.NewAccount)
and
	func NewAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if t := template.PRG(w, r, "account", "NewAccount"); t == nil {
			return // GET template gets the id parameter
		} else if account, err := model.NewAccount(...); err != nil {
			template.Error(...)
		} else {
			t.Set("id", account.ID)
			t.Run() // Redirect
		}
	}

Works with POST, as well as with PUT, and DELETE.
*/
func PRG(w http.ResponseWriter, r *http.Request, dir, base string, inner ...string) (prg *PRGType) {
	switch r.Method {
	case "GET":
		values := r.URL.Query()
		data := make(map[string]interface{}, len(values)+1)
		for key := range values {
			data[key] = r.FormValue(key)
		}
		Run(w, r, dir, base, opt(inner...), data)
	case "PUT", "POST", "DELETE":
		values, _ := url.ParseQuery("")
		prg = &PRGType{&baseType{&values, w, r}}
	}
	return
}
