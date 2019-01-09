package logger

import (
	"os"

	"bitbucket.org/d3dev/parse_pikabu/config"
	logging "github.com/op/go-logging"
)

var Log = logging.MustGetLogger("parse_pikabu/core")
var ParserLog = logging.MustGetLogger("parse_pikabu/parser")
var LogFormat = logging.MustStringFormatter(
	`%{color}%{module} %{pid} %{level:.5s} %{time:15:04:05.000} %{shortfile} %{shortfunc} ▶ %{id:03x}%{color:reset} %{message}`,
)

func init() {
	file, err := os.OpenFile("logs/parse_pikabu.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	// loggingBackend := logger.NewLogBackend(os.Stderr, "", 0)
	loggingBackend := logging.NewLogBackend(file, "", 0)
	loggingBackendFormatter := logging.NewBackendFormatter(loggingBackend, LogFormat)

	logging.SetBackend(loggingBackend, loggingBackendFormatter)

	if config.Settings.Debug {
		logging.SetLevel(logging.DEBUG, "parse_pikabu/core")
	} else {
		logging.SetLevel(logging.INFO, "parse_pikabu/core")
	}
	if config.Settings.Debug {
		logging.SetLevel(logging.DEBUG, "parse_pikabu/parser")
	} else {
		logging.SetLevel(logging.INFO, "parse_pikabu/parser")
	}
}
