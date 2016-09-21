package env

import (
	"os"
)

var envs = make(map[string]string)

func Get(name string) (value string) {
	return envs[name]
}

func set(name, defaultValue string) {
	var env string
	if env = os.Getenv(name); env == "" {
		env = defaultValue
	}
	envs[name] = env
}

func init() {
	set("HTTP_HOST", "localhost")
	set("HTTP_PORT", ":9090")
	set("HTTPS_PORT", ":10443")
	set("DB_HOST", "db1")
	set("DB_PORT", ":28015")
	set("DB_NAME", "expeertise")
}
