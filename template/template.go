package template

import (
	"bytes"
	"github.com/wscherphof/msg"
	"github.com/yosssi/ace"
	"html/template"
	"io"
	"log"
	"net/http"
)

// Run loads and executes a template, writing the output to w
func Run(w io.Writer, r *http.Request, dir, base, inner string, data map[string]interface{}) {
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

// Write loads and executes a template, returning the output as a byte array
func Write(r *http.Request, dir, base, inner string, data map[string]interface{}) []byte {
	var b bytes.Buffer
	Run(&b, r, dir, base, inner, data)
	return b.Bytes()
}
