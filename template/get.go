package template

import (
	"net/http"
)

type GetType struct {
	*baseType
}

// Run executes the template, passing in the set data values.
func (t *GetType) Run(opt_status ...int) {
	response(t.w, t.r, t.dir, t.base, t.inner(), t.data, t.status(opt_status...))
}

/*
GET returns a type that, just as net/url.Values, listens to Set(key, value)
to register string data. Call Run() to execute the template.
*/
func GET(w http.ResponseWriter, r *http.Request, dir, base string, opt_inner ...string) *GetType {
	return &GetType{&baseType{w, r, dir, base, opt_inner, nil}}
}
