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

var Countries = make([]Country, 250)

func LoadCountries () {
  var (
    data []byte
    err error
  )
  if data, err = ioutil.ReadFile("./countries.json"); err == nil {
    if err = json.Unmarshal(data, &Countries); err == nil {
      for i, v := range Countries {
        Countries[i].Code = strings.ToLower(v.Code)
      }
    }
  }
  if err != nil {
    log.Print("ERROR: ", err.Error())
  }
}
