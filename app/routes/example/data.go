package example

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type vk map[string]string

type jsonCountry struct {
	Name string
	Code string
}

var countries vk

func Countries() vk {
	if countries == nil {
		jsonCountries := make([]jsonCountry, 250)
		if data, err := ioutil.ReadFile("/resources/data/example/countries.json"); err != nil {
			log.Panicln("ERROR:", err)
		} else if err := json.Unmarshal(data, &jsonCountries); err != nil {
			log.Panicln("ERROR:", err)
		}
		countries = make(vk, len(jsonCountries))
		for _, v := range jsonCountries {
			v.Code = strings.ToLower(v.Code)
			countries[v.Name] = v.Code
		}
	}
	return countries
}

type jsonTimeZone struct {
	Value string
	Abbr string
	Offset float32
	Isdst bool
	Text string
	Utc []string
}

var timeZones vk

func TimeZones() vk {
	if timeZones == nil {
		jsonTimeZones := make([]jsonTimeZone, 250)
		if data, err := ioutil.ReadFile("/resources/data/example/timezones.json"); err != nil {
			log.Panicln("ERROR:", err)
		} else if err := json.Unmarshal(data, &jsonTimeZones); err != nil {
			log.Panicln("ERROR:", err)
		}
		timeZones = make(vk, len(jsonTimeZones))
		for _, v := range jsonTimeZones {
			if len(v.Utc) > 0 && v.Text != "" {
				timeZones[v.Text] = v.Utc[0]
			}
		}
	}
	return timeZones
}
