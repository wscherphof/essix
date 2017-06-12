package template

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type PRGType struct {
	*BaseType
}

// Run redirects to path/method, including set values.
func (t *PRGType) Run(opt_status ...int) {
	if len(opt_status) == 1 {
		t.Set("_status", strconv.Itoa(opt_status[0]))
	}
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

Works with POST, as well as with PUT, PATCH, and DELETE.
*/
func PRG(w http.ResponseWriter, r *http.Request, dir, base string, opt_inner ...string) (prg *PRGType) {
	if r.Method == "GET" {
		inner := ""
		if len(opt_inner) == 1 {
			inner = opt_inner[0]
		}
		values := r.URL.Query()
		data := make(map[string]interface{}, len(values)+2)
		for key := range values {
			data[key] = r.FormValue(key)
		}
		status := http.StatusOK
		if _status := r.FormValue("_status"); _status != "" {
			status, _ = strconv.Atoi(_status)
		}
		response(w, r, dir, base, inner, data, status)
	} else {
		prg = &PRGType{&BaseType{w, r, dir, base, opt_inner, nil}}
	}
	return
}
