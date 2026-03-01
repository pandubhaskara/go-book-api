package helper

import (
	"log"
	"os"
	"strconv"
)

type logger struct {
	debugLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	env      string
	debug    bool
}

var Logger logger

func init() {
	Logger = NewLogger()
}

func NewLogger() logger {
	flags := log.LstdFlags //log.LstdFlags | log.Lshortfile

	l := logger{
		debugLog: log.New(os.Stdout, "[DEBUG] ", flags),
		infoLog:  log.New(os.Stdout, "[INFO] ", flags),
		warnLog:  log.New(os.Stdout, "[WARNING] ", flags),
		errorLog: log.New(os.Stdout, "[ERROR] ", flags),
		env:      "local",
		debug:    false,
	}

	if val, present := os.LookupEnv("APP_ENV"); present && val != "" {
		l.env = val
	}

	if val, present := os.LookupEnv("APP_DEBUG"); present && val != "" {
		if debug, err := strconv.ParseBool(val); err == nil {
			l.debug = debug
		}
	}

	return l

}

func (l *logger) Debug(v ...interface{}) {
	if l.env != "production" || l.debug {
		l.debugLog.Println(v...)
	}
}

func (l *logger) Debugf(f string, v ...interface{}) {
	if l.env != "production" || l.debug {
		l.debugLog.Printf(f, v...)
	}
}

func (l *logger) Info(v ...interface{}) {
	l.infoLog.Println(v...)
}

func (l *logger) Infof(f string, v ...interface{}) {
	l.infoLog.Printf(f, v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.warnLog.Println(v...)
}

func (l *logger) Warnf(f string, v ...interface{}) {
	l.warnLog.Printf(f, v...)
}

func (l *logger) Error(v ...interface{}) {
	l.errorLog.Println(v...)
}

func (l *logger) Errorf(f string, v ...interface{}) {
	l.infoLog.Printf(f, v...)
}
