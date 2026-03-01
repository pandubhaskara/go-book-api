package helper

import (
	"os"

	"github.com/joho/godotenv"
)

type env struct {
}

func init() {
	l := NewLogger()
	l.Info("Initializing Environment Variables")
	if _, err := os.Stat(".env"); err == nil {
		l.Info("Loading .env file")
		godotenv.Load()
	}
}

func (e *env) Get(k string, d string) (value string) {
	if val, present := os.LookupEnv(k); present && val != "" {
		value = val
	} else {
		value = d
	}
	return
}

var Env env
