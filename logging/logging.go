package logging

import (
	"github.com/op/go-logging"
	"os"
)

var Log = logging.MustGetLogger("parse_pikabu")
var logFormat = logging.MustStringFormatter(
	`%{color}%{level:.5s} %{time:15:04:05.000} %{shortfunc} ▶ %{id:03x}%{color:reset} %{message}`,
)

func init() {
	file, err := os.OpenFile("logs/main.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	// loggingBackend := logging.NewLogBackend(os.Stderr, "", 0)
	loggingBackend := logging.NewLogBackend(file, "", 0)
	loggingBackendFormatter := logging.NewBackendFormatter(loggingBackend, logFormat)

	logging.SetBackend(loggingBackend, loggingBackendFormatter)
	Log.Debug("app started")
}
