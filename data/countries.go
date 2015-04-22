package data

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

var countries []Country

func Countries () *[]Country {
  if len(countries) == 0 {
    countries = make([]Country, 250)
    if data, err := ioutil.ReadFile("./data/countries.json"); err != nil {
      log.Panicln("ERROR:", err)
    } else if err := json.Unmarshal(data, &countries); err != nil {
      log.Panicln("ERROR:", err)
    }
    // Lowercase all the codes!
    for i, v := range countries {
      countries[i].Code = strings.ToLower(v.Code)
    }
  }
  return &countries
}
