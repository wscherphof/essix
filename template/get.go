package template

import (
	"net/http"
)

type getType struct {
	*prgType
	dir   string
	base  string
	inner string
}

func (t *getType) Run() {
	data := make(map[string]interface{})
	for key := range *t.Values {
		data[key] = t.Values.Get(key)
	}
	Run(t.w, t.r, t.dir, t.base, t.inner, data)
}

func GET(w http.ResponseWriter, r *http.Request, dir, base string, opt_inner ...string) *getType {
	values := r.URL.Query()
	return &getType{
		prgType: &prgType{
			Values: &values,
			w:      w,
			r:      r,
		},
		dir:   dir,
		base:  base,
		inner: inner(opt_inner...),
	}
}
