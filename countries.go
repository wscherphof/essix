package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "strings"
)

type Country struct {
  Name string
  Code string
}

var _countries []Country

func Countries () *[]Country {
  if _countries == nil {
    _countries = make([]Country, 250)
    data, err := ioutil.ReadFile("./countries.json")
    if err != nil {
      log.Panicln("ERROR:", err.Error())
    }
    err = json.Unmarshal(data, &_countries)
    if err != nil {
      log.Panicln("ERROR:", err.Error())
    }
    for i, v := range _countries {
      _countries[i].Code = strings.ToLower(v.Code)
    }
  }
  return &_countries
}
