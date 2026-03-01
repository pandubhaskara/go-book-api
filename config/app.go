package config

import (
	"strconv"

	"go-book-api/helper"
)

var App = struct {
	Env      string
	Name     string
	Ver      string
	Key      string
	Debug    bool
	Timezone string
	Host     string
	Port     string
	Cors     bool
}{
	Env:      helper.Env.Get("APP_ENV", "local"),
	Name:     helper.Env.Get("APP_NAME", ""),
	Ver:      helper.Env.Get("APP_VER", ""),
	Key:      helper.Env.Get("APP_KEY", ""),
	Debug:    false,
	Timezone: helper.Env.Get("APP_TIMEZONE", "Etc/UTC"),
	Host:     helper.Env.Get("APP_HOST", "localhost"),
	Port:     helper.Env.Get("APP_PORT", "8080"),
	Cors:     false,
}

func init() {
	if debug, err := strconv.ParseBool(helper.Env.Get("APP_DEBUG", "false")); err == nil {
		App.Debug = debug
	}

	if cors, err := strconv.ParseBool(helper.Env.Get("APP_CORS", "false")); err == nil {
		App.Cors = cors
	}
}
