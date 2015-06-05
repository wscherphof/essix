package util

import (
	"bytes"
	"github.com/wscherphof/msg"
	"github.com/yosssi/ace"
	"html/template"
	"io"
	"log"
	"net/http"
)

func aceOptions(r *http.Request) *ace.Options {
	return &ace.Options{
		BaseDir: "templates",
		FuncMap: template.FuncMap{
			"Msg": msg.Msg(r),
		},
	}
}

func Template(base string, inner string, data map[string]interface{}) func(io.Writer, *http.Request) {
	if data == nil {
		data = map[string]interface{}{}
	}
	return func(w io.Writer, r *http.Request) {
		lang := msg.Language(r)
		data["lang"] = lang
		if inner == "lang" {
			inner = base + "-" + lang.Main
		}
		if tpl, err := ace.Load(base, inner, aceOptions(r)); err != nil {
			log.Panicln("ERROR: ace.Load:", err)
		} else if err := tpl.Execute(w, data); err != nil {
			log.Panicln("ERROR: tpl.Execute:", err)
		}
	}
}

func BTemplate(base string, inner string, data map[string]interface{}) func(*http.Request) []byte {
	var b bytes.Buffer
	return func(r *http.Request) []byte {
		Template(base, inner, data)(&b, r)
		return b.Bytes()
	}
}
