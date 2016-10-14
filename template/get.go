package template

import (
	"net/http"
)

type getType struct {
	*baseType
	dir   string
	base  string
	inner string
}

func (t *getType) Run() {
	data := make(map[string]interface{}, len(*(t.Values)) + 1)
	for key := range *t.Values {
		data[key] = t.Values.Get(key)
	}
	Run(t.w, t.r, t.dir, t.base, t.inner, data)
}

func GET(w http.ResponseWriter, r *http.Request, dir, base string, inner ...string) *getType {
	values := r.URL.Query()
	return &getType{
		baseType: &baseType{&values, w, r},
		dir:      dir,
		base:     base,
		inner:    opt(inner...),
	}
}
