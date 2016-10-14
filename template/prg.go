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

type prgType struct {
	*baseType
}

func (t *prgType) Run() {
	path := t.r.URL.Path
	path += "/" + strings.ToLower(t.r.Method)
	path += "?" + t.Values.Encode()
	http.Redirect(t.w, t.r, path, http.StatusSeeOther)
}

func PRG(w http.ResponseWriter, r *http.Request, dir, base string, inner ...string) (prg *prgType) {
	switch r.Method {
	case "GET":
		values := r.URL.Query()
		data := make(map[string]interface{}, len(values) + 1)
		for key := range values {
			data[key] = r.FormValue(key)
		}
		Run(w, r, dir, base, opt(inner...), data)
	case "PUT", "POST", "DELETE":
		values, _ := url.ParseQuery("")
		prg = &prgType{&baseType{&values, w, r}}
	}
	return
}
