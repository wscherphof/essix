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
	data := make(map[string]interface{}, len(*(t.Values))+1)
	for key := range *t.Values {
		data[key] = t.Values.Get(key)
	}
	Run(t.w, t.r, t.dir, t.base, t.inner, data)
}

/*
GET returns a type that, just as net/url.Values, listens to Set(key, value)
to register string data. Call Run() to execute the template.
*/
func GET(w http.ResponseWriter, r *http.Request, dir, base string, inner ...string) *GetType {
	values := r.URL.Query()
	return &GetType{
		baseType: &baseType{&values, w, r},
		dir:      dir,
		base:     base,
		inner:    opt(inner...),
	}
}
