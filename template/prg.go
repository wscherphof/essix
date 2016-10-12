package template

import (
	"net/http"
	"net/url"
	"strings"
)

type prgType struct {
	*url.Values
	w http.ResponseWriter
	r *http.Request
}

func (p *prgType) Redirect() {
	path := p.r.URL.Path
	path += "/" + strings.ToLower(p.r.Method)
	path += "?" + p.Values.Encode()
	http.Redirect(p.w, p.r, path, http.StatusSeeOther)
}

func PRG(w http.ResponseWriter, r *http.Request, dir, base, inner string) (prg *prgType) {
	switch r.Method {
	case "GET":
		data := make(map[string]interface{})
		for key, _ := range r.URL.Query() {
			data[key] = r.FormValue(key)
		}
		Run(w, r, dir, base, inner, data)
	case "PUT", "POST", "DELETE":
		values, _ := url.ParseQuery("")
		prg = &prgType{
			Values: &values,
			w: w,
			r: r,
		}
	}
	return
}
