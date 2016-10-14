/*
Package env manages environment variables.
*/
package env

import (
	"log"
	"os"
	"strconv"
)

/*
Get returns the value of the named environment variable,
or sets and returns the given default value,
or else logs a fatal error.
Use "" as default value to not set the environment variable if missing.
*/
func Get(name string, defaultValue ...string) (value string) {
	if value = os.Getenv(name); value == "" {
		if len(defaultValue) == 1 {
			value = defaultValue[0]
			if value != "" {
				Set(name, value)
			}
		} else {
			log.Fatal("ERROR: Environment variable ", name, " not set")
		}
	}
	return
}

/*
Set sets an environment variable.
*/
func Set(name, value string) {
	os.Setenv(name, value)
}

/*
GetInt is like Get, but with integer values. A fatal error is logged if the
variable was set with a non-integer value.
*/
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
