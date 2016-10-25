/*
Package template renders templates using github.com/yossi/ace, providing text
translation through github.com/wscherphof/msg.

It provides an httprouter.Handle to just render a template as the action for a route.

It renders and sends email messages.

It renders handled responses to GET requests.

It renders errors from request handlers.

It implements the Post-Redirect-Get pattern, redirecting to a GET request,
rendering the template with data from a POST, PUT, or DELETE request
*/
package template

import (
	"bytes"
	"github.com/wscherphof/msg"
	"github.com/wscherphof/secure"
	"github.com/yosssi/ace"
	"io"
	"log"
	"net/http"
)

const location = "/resources/templates"

/*
String loads and executes a template, returning the output.

dir is a directory (one deep) under /resources/templates

base is the template name (filename without extension) in dir

inner is the inner template name (without extension) in dir

Both base and inner may include paths relative to dir.

data is the pipeline data to pass to the template
*/
func String(r *http.Request, dir, base, inner string, opt_data ...map[string]interface{}) string {
	var data map[string]interface{}
	if len(opt_data) == 1 {
		data = opt_data[0]
	}
	var b bytes.Buffer
	run(&b, r, dir, base, inner, data)
	return string(b.Bytes())
}

func response(w http.ResponseWriter, r *http.Request, dir, base, inner string, data map[string]interface{}, opt_status ...int) {
	status := http.StatusOK
	if len(opt_status) == 1 {
		status = opt_status[0]
	}
	if w.Header().Get("Content-Type") == "" {
		// Set the Content-Type to prevent CompressHandler from doing so after our WriteHeader()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	w.WriteHeader(status)
	run(w, r, dir, base, inner, data)
}

func run(w io.Writer, r *http.Request, dir, base, inner string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{}, 2)
	}
	translator := msg.Translator(r)
	data["msg"] = translator
	if r.Method == "GET" && data["_formtoken"] == nil {
		formtoken := secure.NewFormToken(r)
		data["_formtoken"] = formtoken.String()
	}
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
