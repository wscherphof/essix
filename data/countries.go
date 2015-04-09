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

var _countries []Country

func Countries () *[]Country {
  if len(_countries) == 0 {
    _countries = make([]Country, 250)
    if data, err := ioutil.ReadFile("./data/countries.json"); err != nil {
      log.Panicln("ERROR:", err)
    } else if err := json.Unmarshal(data, &_countries); err != nil {
      log.Panicln("ERROR:", err)
    }
    // Lowercase all the codes!
    for i, v := range _countries {
      _countries[i].Code = strings.ToLower(v.Code)
    }
  }
  return &_countries
}
