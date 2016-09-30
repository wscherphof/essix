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

func Template(w io.Writer, r *http.Request, dir, base, inner string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	lang := msg.Language(r)
	if inner == "lang" {
		inner = base + "-" + lang.Main
	}
	data["lang"] = lang
	var options = &ace.Options{
		BaseDir: "/resources/templates/" + dir,
		FuncMap: template.FuncMap{
			"Msg": msg.Msg(r),
		},
	}
	if template, err := ace.Load(base, inner, options); err != nil {
		log.Panicln("ERROR: ace.Load:", err)
	} else if err := template.Execute(w, data); err != nil {
		log.Panicln("ERROR: template.Execute:", err)
	}
}

func BTemplate(r *http.Request, dir, base, inner string, data map[string]interface{}) []byte {
	var b bytes.Buffer
	Template(&b, r, dir, base, inner, data)
	return b.Bytes()
}
