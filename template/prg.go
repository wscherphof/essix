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

func inner(opt_inner ...string) (ret string) {
	if len(opt_inner) == 1 {
		ret = opt_inner[0]
	}
	return
}

func PRG(w http.ResponseWriter, r *http.Request, dir, base string, opt_inner ...string) (prg *prgType) {
	switch r.Method {
	case "GET":
		data := make(map[string]interface{})
		for key := range r.URL.Query() {
			data[key] = r.FormValue(key)
		}
		Run(w, r, dir, base, inner(opt_inner...), data)
	case "PUT", "POST", "DELETE":
		values, _ := url.ParseQuery("")
		prg = &prgType{&baseType{&values, w, r}}
	}
	return
}
