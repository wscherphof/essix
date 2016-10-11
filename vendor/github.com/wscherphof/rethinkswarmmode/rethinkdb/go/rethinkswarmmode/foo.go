package main

import (
	"fmt"
	r "gopkg.in/dancannon/gorethink.v2"
	"html"
	"log"
	"net/http"
)

func main() {

	var url = "db1:28015"

	session, err := r.Connect(r.ConnectOpts{
		Address: url,
	})
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/bar", func(w http.ResponseWriter, req *http.Request) {

		res, err := r.Expr("Hello from Rethink").Run(session)
		if err != nil {
			log.Fatalln(err)
		}

		var response string
		err = res.One(&response)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Fprintf(w, "Hello, %q 0.1\n", html.EscapeString(req.URL.Path))
		fmt.Fprintf(w, response+"\n")
	})

	log.Fatal(http.ListenAndServe(":9090", nil))
}
