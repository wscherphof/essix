package template

import (
	"net/http"
)

type BaseType struct {
	w         http.ResponseWriter
	r         *http.Request
	dir       string
	base      string
	opt_inner []string
	data      map[string]interface{}
}

/*
Set saves data to the template's pipeline.
*/
func (b *BaseType) Set(key string, value interface{}) {
	if b.data == nil {
		b.data = make(map[string]interface{})
	}
	b.data[key] = value
}

func (b *BaseType) inner() string {
	if len(b.opt_inner) == 1 {
		return b.opt_inner[0]
	}
	return ""
}
