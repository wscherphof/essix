package example

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

/*
A (value, key) map. In that order, so the template will sort on the values when
ranging.
*/
type vk map[string]string

var countries vk

type jsonCountry struct {
	Name string
	Code string
}

func Countries() vk {
	if countries == nil {
		var data []jsonCountry
		readJsonData(&data, "/resources/data/example/countries.json")
		countries = make(vk, len(data))
		for _, v := range data {
			v.Code = strings.ToLower(v.Code)
			countries[v.Name] = v.Code
		}
	}
	return countries
}

var timeZones vk

type jsonTimeZone struct {
	Value  string
	Abbr   string
	Offset float32
	Isdst  bool
	Text   string
	Utc    []string
}

func TimeZones() vk {
	if timeZones == nil {
		var data []jsonTimeZone
		readJsonData(&data, "/resources/data/example/timezones.json")
		timeZones = make(vk, len(data))
		for _, v := range data {
			if len(v.Utc) > 0 && v.Text != "" {
				timeZones[v.Text] = v.Utc[0]
			}
		}
	}
	return timeZones
}

func readJsonData(dst interface{}, file string) {
	if data, err := ioutil.ReadFile(file); err != nil {
		log.Panicln("ERROR:", err)
	} else if err := json.Unmarshal(data, dst); err != nil {
		log.Panicln("ERROR:", err)
	}
}
