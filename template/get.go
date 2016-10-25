package template

import (
	"net/http"
)

type GetType struct {
	*baseType
	dir   string
	base  string
	inner string
}

// Run executes the template, passing in the set data values.
func (t *GetType) Run() {
	run(t.w, t.r, t.dir, t.base, t.inner, t.data, t.Status)
}

/*
GET returns a type that, just as net/url.Values, listens to Set(key, value)
to register string data. Call Run() to execute the template.
*/
func GET(w http.ResponseWriter, r *http.Request, dir, base string, inner ...string) *GetType {
	return &GetType{
		baseType: newBaseType(w, r),
		dir:      dir,
		base:     base,
		inner:    opt(inner...),
	}
}
