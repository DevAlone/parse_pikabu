package logger

import (
	"os"

	"bitbucket.org/d3dev/parse_pikabu/core/config"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger
var logFile *os.File

var PikagoLog *logrus.Logger
var pikagoLogFile *os.File

var PikagoHttpLog *logrus.Logger
var pikagoHttpLogFile *os.File

func Init() {
	var err error
	logFile, err = os.OpenFile("logs/parser.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	pikagoLogFile, err = os.OpenFile("logs/pikago.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	pikagoHttpLogFile, err = os.OpenFile("logs/http.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	Log = logrus.New()
	PikagoLog = logrus.New()
	PikagoHttpLog = logrus.New()

	for log, file := range map[*logrus.Logger]*os.File{
		Log:           logFile,
		PikagoLog:     pikagoLogFile,
		PikagoHttpLog: pikagoHttpLogFile,
	} {
		log.SetOutput(file)
		// log.SetFormatter(&logrus.JSONFormatter{})
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
		if config.Settings.Debug {
			log.SetLevel(logrus.DebugLevel)
		} else {
			log.SetLevel(logrus.WarnLevel)
		}
	}
}

func Cleanup() {
	_ = logFile.Close()
}
