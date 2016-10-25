package template

import (
	"net/http"
)

type baseType struct {
	w         http.ResponseWriter
	r         *http.Request
	dir       string
	base      string
	opt_inner []string
	data      map[string]interface{}
}

func (b *baseType) Set(key string, value interface{}) {
	if b.data == nil {
		b.data = make(map[string]interface{})
	}
	b.data[key] = value
}

func (b *baseType) inner() string {
	if len(b.opt_inner) == 1 {
		return b.opt_inner[0]
	}
	return ""
}

func (b *baseType) status(opt_status ...int) int {
	if len(opt_status) == 1 {
		return opt_status[0]
	}
	return http.StatusOK
}
