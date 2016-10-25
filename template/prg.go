package template

import (
	"net/http"
	"net/url"
	"strings"
)

type PRGType struct {
	*baseType
}

// Run redirects to path/method, including set values.
func (t *PRGType) Run(opt_status ...int) {
	t.data["_status"] = t.status(opt_status...)
	path := t.r.URL.Path
	path += "/" + strings.ToLower(t.r.Method)
	path += "?"
	for k, v := range t.data {
		s := v.(string)
		path += k + "=" + url.QueryEscape(s) + "&"
	}
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
func PRG(w http.ResponseWriter, r *http.Request, dir, base string, opt_inner ...string) (prg *PRGType) {
	prg = &PRGType{&baseType{w, r, dir, base, opt_inner, nil}}
	if r.Method == "GET" {
		values := r.URL.Query()
		data := make(map[string]interface{}, len(values)+2)
		for key := range values {
			data[key] = r.FormValue(key)
		}
		response(w, r, dir, base, prg.inner(), data, data["_status"].(int))
	}
	return
}
