package template

import (
	"bytes"
	"github.com/wscherphof/msg"
	"github.com/yosssi/ace"
	"io"
	"log"
	"net/http"
)

const location = "/resources/templates"

// Run loads and executes a template, writing the output to w
func Run(w io.Writer, r *http.Request, dir, base, inner string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	translator := msg.Translator(r)
	data["msg"] = translator
	var options = &ace.Options{
		BaseDir: location + "/" + dir,
	}
	if inner == "lang" {
		if file, err := translator.File(location, dir, base); err != nil {
			log.Panicf("ERROR: no lang template found for %s/%s/%s %+v", location, dir, base, err)
			return
		} else {
			inner = file
		}
	}
	if template, err := ace.Load(base, inner, options); err != nil {
		log.Panicln("ERROR: ace.Load:", err)
	} else if err := template.Execute(w, data); err != nil {
		log.Panicln("ERROR: template.Execute:", err)
	}
}

// Write loads and executes a template, returning the output
func Write(r *http.Request, dir, base, inner string, data map[string]interface{}) string {
	var b bytes.Buffer
	Run(&b, r, dir, base, inner, data)
	return string(b.Bytes())
}
