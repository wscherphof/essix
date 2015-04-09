package main

import (
  "net/http"
  "net/http/httptest"
  "log"
  "testing"
  "io/ioutil"
  "github.com/wscherphof/msg"
)

var m, a = msg.Init()
type template_test_input struct{
  accept_language string
  want string
}

func TestLanguage(t *testing.T) {
  template_test_run(t, "template_test_language", []template_test_input{
    template_test_input{accept_language: "nl-nl", want: "nl"},
    template_test_input{accept_language: "en-gb", want: "en"},
  })
}

func TestMsg(t *testing.T) {
  m("hello")
  a("nl", "hallo")
  a("en", "hello")
  template_test_run(t, "template_test_msg", []template_test_input{
    template_test_input{accept_language: "nl-nl", want: "hallo"},
    template_test_input{accept_language: "en-gb", want: "hello"},
  })
}

func template_test_run(t *testing.T, template string, inputs []template_test_input) {
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    T(template, "", nil)(w, r, nil)
  }))
  defer ts.Close()
  client := &http.Client{}
  req, err := http.NewRequest("GET", ts.URL, nil)
  if err != nil {
    log.Fatal(err)
  }
  for _, v := range inputs {
    req.Header.Set("Accept-Language", v.accept_language)
    res, err := client.Do(req)
    if err != nil {
      log.Fatal(err)
    }
    content, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
      log.Fatal(err)
    }
    got := string(content)
    if got != v.want {
      t.Error("accept_language:", v.accept_language, "want:", v.want, "got:", got)
    }
  }
}
