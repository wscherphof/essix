package env

import (
	"os"
)

func Get(name string) (value string) {
	return os.Getenv(name)
}

func Set(name, value string) {
	os.Setenv(name, value)
}

func Default(name string, defaultValue string) (value string) {
	var env := Get(name)
	if env == "" {
		env = defaultValue
		Set(name, env)
	}
	return env
}
