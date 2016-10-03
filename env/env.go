package env

import (
	"os"
	"log"
	"strconv"
)

func Get(name string, defaultValue ...string) (value string) {
	if value = os.Getenv(name); value == "" {
		if len(defaultValue) == 1 {
			value = defaultValue[0]
		} else {
			log.Fatal("ERROR: Environment variable " , name, " not set")
		}
	}
	return
}

func Set(name, value string) {
	os.Setenv(name, value)
}

func GetInt(name string, defaultValue ...int) (value int) {
	var stringVal string
	if len(defaultValue) == 1 {
		stringVal = Get(name, strconv.Itoa(defaultValue[0]))
	} else {
		stringVal = Get(name)
	}
	var err error
	if value, err = strconv.Atoi(stringVal); err != nil {
		log.Fatal("ERROR: Environment variable ", name, " has value \"", stringVal, "\" while an integer value is expected")
	}
	return 
}
